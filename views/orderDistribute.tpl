<!DOCTYPE html>
<html lang="en">
<head>
<title>抢订单</title>
<script src="javascripts/jquery.js"></script>
<script src="javascripts/underscore.js"></script>
<script type="text/javascript" src="http://api.map.baidu.com/api?type=quick&ak=kU4NWwyP5SwguC2W2WAfO1bO&v=1.0"></script>

<style type="text/css">

</style>
</head>
<body>
    <h1>{{.distributor.Name}}</h1>
    <div>
        <input id="btnPrepared" type="button" value="准备完毕" onclick="prepared()" style="margin-bottom: 10px;"></br>
        <input id="btnAccept" type="button" value="接受订单" onclick="selectOrder()" style="margin-bottom: 10px;">
        <input id="btnPass" type="button" value="暂不考虑" onclick="" style="margin-bottom: 10px;margin-left: 70px;"></br>
        <input id="btnStartDistribute" type="button" value="开始配送" onclick="startDistribution()" style="margin-bottom: 10px"></br>

        <div>
            <select id="orderIDList" style="width: 400px; font-size: 30px;">
<!--             <option value="volvo">Volvo</option>
            <option value="saab">Saab</option>
            <option value="fiat">Fiat</option>
            <option value="audi">Audi</option>
 -->            
            </select>
        </div>
        <div id="allmap" style="height:300px;width:500px;margin-top:10px;"></div>
    </div>
<!--     <div id="btnSelectOrder" onclick="selectOrder()" style="text-decoration: underline; padding-bottom: 10px; font-weight: 500; font-size: 25px; color: white;">
      抢订单
    </div>
 -->
    <div style="  font-size: 23px; margin-bottom: 10px; margin-top: 20px;">信息提示</div>
     <div id="output"></div>

</body>
<script type="text/javascript">
    var conn;
    var output;
    var distributorID = "{{.distributor.ID}}"
    var map = new BMap.Map("allmap");
    var marker = null
    var orders = []

    $(function() {
        output = $("#output")
        hideAllControls()
        appendLog("正在初始化基础数据")
        $.get("/distributors?id={{.distributor.ID}}",function(data){
            console.log(data)
            if(_.size(data) > 0){
                var distributor = data[0]
                console.log(distributor)
                if(distributor.CheckPoint <= 0){//还在初始化阶段
                    $("#btnPrepared").show()

                }else{//已经初始化过，中间可能掉线
                    showOrderSelectButton()
                }
            }
        })

        $("#orderIDList").change(function(){//code...
            var orderID = $("#orderIDList").find("option:selected").val();
            if(orderID == null || orderID.length <= 0){
                // alert("订单不能为空")
                return
            }
            var order = _.findWhere(orders, {ID: orderID})
            if(order != null){
                var pos = order.GeoSrc
                if(pos != null){
                    setMapMarker(pos.Lng, pos.Lat, true)
                }else{
                    setMapMarker(0, 0)
                }
                // if(pos != null && pos.Lat.length > 0 && pos.Lng.length > 0){
                //     setMapMarker(parseFloat(pos.Lng),parseFloat(pos.Lat),true)
                // }else{
                //     setMapMarker(0, 0)
                // }
            }
        });
        prepareConn()

        // var point = new BMap.Point(); 
        setMapMarker(116.404, 39.915, false)
                    // $("#btnSelectOrder").css("color","black")
        // var selectOptions = $("#orderIDList option").remove()
        // var selectOptionCount = selectOptions.length
        // console.log("%d 个 option", selectOptionCount)
        // for (var i = 0; i < selectOptionCount; i++) {
        //     console.log(selectOptions[i].text)
        // };

    })
    function setMapMarker(lng,lat, bAddMarker){
        map.removeOverlay(marker)
        if(lng > 0 || lat > 0){
            map.centerAndZoom(new BMap.Point(lng, lat), 12);
            map.addControl(new BMap.ZoomControl());  //添加地图缩放控件
            if(bAddMarker == true){
                marker = new BMap.Marker(new BMap.Point(lng, lat));  //创建标注
                map.addOverlay(marker);                 // 将标注添加到地图中            
            }
        }
    }

    function appendLog(message) {
        // if(output.children().length > 20){
        //     output.empty()
        // }
        output.prepend(message+"</br>");
    }    
    function prepareConn(){
        if (window["WebSocket"]) {
            conn = new WebSocket("ws://{{.HOST}}/wsOrderDistribution?id={{.distributor.ID}}");
            conn.onclose = function(evt) {
                appendLog("与服务器连接连接关闭，刷新重试")
            }
            conn.onopen = function(evt){
                appendLog("与服务器连接成功")
            }
            conn.onmessage = function(evt) {
                // appendLog($("<div/>").text(evt.data))
                var msg = evt.data
                console.log(msg)
                msg = JSON.parse(msg)
                switch(msg.MessageType){
                    case {{.pro_order_distribution_proposal}}://订单分发
                    var currentIndex = $("#orderIDList").get(0).selectedIndex
                    console.log("新订单推送到，当前选择的订单的索引为 %d", currentIndex)
                    orders = msg.Data
                    $("#orderIDList option").remove()
                    _.each(orders, function(order){
                        var orderTip = "编号："+order.ID
                        if (order.GeoSrc != null){
                            orderTip += "  位置:"+order.GeoSrc.Address
                        } 
                        $("#orderIDList").append('<option value="'+ order.ID +'">'+ orderTip +'</option>');
                    })
                    if(currentIndex < 0){
                        currentIndex = 0
                    }
                    $("#orderIDList").get(0).selectedIndex = currentIndex
                    $("#orderIDList").trigger("change")
                    break
                    case {{.pro_timer_count_down}}://计时
                    appendLog("-> "+ msg.Data)
                    break
                    case {{.pro_begin_to_select_order}}://开始选择订单
                    appendLog("开始！")
                    $("#btnSelectOrder").css("color","black")
                    break
                    case {{.pro_message_broadcast}}://消息广播
                    appendLog(msg.Data)
                    break
                    case {{.pro_order_select_result}}://订单分配结果
                    if(msg.Data.DistributorID == distributorID){
                        appendLog("抢到了订单 "+msg.Data.OrderID)
                    }else{
                        appendLog("没有抢到订单 "+msg.Data.OrderID)
                    }
                    console.log(msg.Data)
                    break
                    case {{.pro_distribution_prepared}}://订单满载，可以准备配送了
                    if(msg.Data == distributorID){
                        hideOrderSelectButton()
                        $("#btnStartDistribute").show()
                    }

                    break
                }
                // distributionProposals = JSON.parse(distributionProposals)
            }
        } else {
            appendLog("浏览器不支持")
        }        
    }
    function send(msg){
        if (!conn) {
            return false;
        }
        conn.send(JSON.stringify(msg))
    }
    function selectOrder(){
        orderID = $("#orderIDList").find("option:selected").val();
        if(orderID == null || orderID.length <= 0){
            alert("订单不能为空")
            return
        }
        var msg = {MessageType: {{.pro_order_select_response}}, Data:{OrderID: orderID, DistributorID: distributorID}}
        send(msg)
        // $("#btnSelectOrder").css("color","white")
    }
    function prepared(){
        var msg = {MessageType: {{.pro_prepared_for_select_order}}, Data:{DistributorID: distributorID}}
        send(msg)
        prepareSelectOrderControls()
    }
    function prepareSelectOrderControls(){
        $("#btnPrepared").hide()
        showOrderSelectButton()
    }
    function hideOrderSelectButton(){
        $("#btnAccept").hide()
        $("#btnPass").hide()
        $("#orderIDList").hide()
        $("#btnStartDistribute").hide()
    }
    function hideAllControls(){
        $("#btnPrepared").hide()
        $("#btnAccept").hide()
        $("#btnPass").hide()
        $("#orderIDList").hide()
        $("#btnStartDistribute").hide()
        $("#allmap").hide()
    }
    function showOrderSelectButton(){
        $("#btnAccept").show()
        $("#btnPass").show()
        $("#orderIDList").show()
        $("#allmap").show()
    }
    function startDistribution(){
        window.location.href = "/distribution?id={{.distributor.ID}}"
    }
</script>

</html>
