<!DOCTYPE html>
<html lang="en">
<head>
<title>配送订单</title>
<script src="javascripts/jquery.js"></script>
    <script src="javascripts/lodash.js"></script>
<!-- <script src="javascripts/underscore.js"></script> -->
<script type="text/javascript" src="http://api.map.baidu.com/api?type=quick&ak=kU4NWwyP5SwguC2W2WAfO1bO&v=1.0"></script>

<style type="text/css">

</style>
</head>
<body>
    <h1>{{.distributor.Name}}</h1>
    <div>
        <!-- <input id="btnPrepared" type="button" value="准备完毕" onclick="prepared()" style="margin-bottom: 10px;"></br> -->
        <!-- <input id="btnAccept" type="button" value="接受订单" onclick="selectOrder()" style="margin-bottom: 10px;"> -->
        <!-- <input id="btnPass" type="button" value="暂不考虑" onclick="" style="margin-bottom: 10px;margin-left: 70px;"></br> -->
        <div>
            <div style="margin-bottom:5px;font-size:18px;font-weight: 500;"> 当前位置： </div>
            <span id= "currentLng">0</span>
            <span id= "currentLat" style="margin-left: 20px;">0</span>

            <div style="margin-bottom:5px;font-size:18px;font-weight: 500;margin-top:15px;"> 目标位置： </div>
            <span id= "destLng">0</span>
            <span id= "destLat" style="margin-left: 20px;">0</span>
            
            <div style="margin-bottom:5px;font-size:18px;font-weight: 500;margin-top:15px;"> 最高时速： </div>
            <span id= "maxSpeed" style="margin-left: 0px;">0</span>
            
            <div style="margin-bottom:5px;font-size:18px;font-weight: 500;margin-top:15px;"> 当前时速： </div>
            <span id= "currentSpeed" style="margin-left: 0px;">0</span>

            <div style="margin-bottom:5px;font-size:18px;font-weight: 500;margin-top:15px;"> 可选路径节点： </div>
            <select multiple id="nodeSelect" style="width:60%;"></select></br>
            <input id="btnChangeMoveState" type="button" value="启动" state = "0" onclick="onChangeMoveState()" style="margin-bottom: 10px; width: 100px; height: 30px; font-size: 20px; margin-top: 3px;"></br>


        </div>
        <div>
        </div>
        <div id="allmap" style="height:200px;width:500px;margin-top:10px;"></div>
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
    var marker = null
    var distributor = null
    var mapData = null
    var currentLng,currentLat, maxSpeed, currentSpeed, destLng, destLat
    var map = new BMap.Map("allmap");

    $(function() {
        output = $("#output")
        currentLng = $("#currentLng")
        currentLat = $("#currentLat")
        destLng = $("#destLng")
        destLat = $("#destLat")

        maxSpeed = $("#maxSpeed")
        currentSpeed = $("#currentSpeed")
        // mapInit()
        // hideAllControls()
        appendLog("正在初始化基础数据")
        prepareConn()
        initMap()
        resetMap2Initial()

        $("#nodeSelect").dblclick(function(){
            var option = $("option:selected", this)
            var pointID = option.val()
            console.log("选中路径点 ",pointID)
            //将目标点设置为选中的点，发送给服务器
            var p = _.find(mapData.Points, {ID: parseInt(pointID)})
            if(p == null){
                console.warn("没有查找到选中的点")
                return
            }else{
                // dCloned = _.clone(distributor, true)
                // dCloned.DestPos = p
                send({MessageType: {{.pro_reset_destination_request}}, Data: {PositionID: p.ID, DistributorID: distributor.ID}})//请求重置目标点
            }
        })
    })
    function onwsMessage(evt){
        var msg = evt.data
        console.log(msg)
        msg = JSON.parse(msg)
        switch(msg.MessageType){
            case {{.pro_distribution_prepared}}://pro_distribution_prepared 地图基础数据
            console.dir(msg.Data)
            distributor = msg.Data.Distributor
            mapData = msg.Data.MapData
            setDistributorInfoShow()
            refreshNodeToSelect()
            break
            case {{.pro_reset_destination}}://pro__reset_destination
            console.info("目标点重置")
            distributor = msg.Data
            setDistributorInfoShow()
            refreshNodeToSelect()
            break
            case {{.pro_change_state}}://pro_reset_change_state
            distributor = msg.Data
            console.log("currentSpeed: ", distributor.CurrentSpeed)
            setDistributorInfoShow()
            refreshNodeToSelect()
            break
            case {{.pro_move_to_new_position}}:
            distributor = msg.Data
            console.log("move to new possition :", distributor.CurrentPos)
            setDistributorInfoShow()
            break
        }
        // distributionProposals = JSON.parse(distributionProposals)        
    }
    function onChangeMoveState(){
        var btnChangeMoveState = $("#btnChangeMoveState")
        var state = btnChangeMoveState.attr("state")
        if(state == "0"){//未运行状态
            send({MessageType: {{.pro_change_state_request}}, Data: {State: 1, DistributorID: distributor.ID}})//请求重置运动状态
        }else{
            send({MessageType: {{.pro_change_state_request}}, Data: {State: 0, DistributorID: distributor.ID}})//请求重置运动状态
        }
    }
    function setDistributorInfoShow(){
        //标识当前位置
        if(distributor.CurrentPos != null){
            var pos = distributor.CurrentPos
            currentLng.text(pos.Lng.toFixed(6))
            currentLat.text(pos.Lat.toFixed(6))
            setMapMarker(pos.Lng, pos.Lat, true)
        }else{
            resetMap2Initial()
            console.error("当前位置不能为空")
            return
        }
        if(distributor.DestPos != null){
            var pos = distributor.DestPos
            destLng.text(pos.Lng)
            destLat.text(pos.Lat)
        }else{
            destLng.text(0)
            destLat.text(0)
        }
        maxSpeed.text(distributor.Speed)
        // currentSpeed.text(distributor.CurrentSpeed)
        var btnChangeMoveState = $("#btnChangeMoveState")
        if(distributor.CurrentSpeed <= 0){
            currentSpeed.text(distributor.CurrentSpeed)
            btnChangeMoveState.attr("state", "0")
            btnChangeMoveState.val("启动")
        }else{
            currentSpeed.text(distributor.CurrentSpeed)
            btnChangeMoveState.attr("state", "1")
            btnChangeMoveState.val("停止")
        }      
    }
    function refreshNodeToSelect(){

        //查找可以走向的路径节点
        //这里有两种情况，正处于路径节点上和在两个节点之间
        //对于第一种情况，应该查找所有与该点相关的路径
        //对于第二种情况，显示所在路径的起点与终点
        if(isDistributorOnNode(distributor) == true){
            var pos = distributor.CurrentPos
            var lines = _.filter(mapData.Lines, function(l){
                var r = _.find([l.Start, l.End], function(p){
                    return p.Lat == pos.Lat && p.Lng == pos.Lng
                })
                return r != null
            })
            console.log("filter lines :", lines)
            var points = _.map(lines, function(l){
                return [l.Start, l.End]
            })
            // points = _.flatten(points)
            // points = _.filter(points, )
            points = _.chain(points).flatten().filter(function(p){
                return p.Lat != pos.Lat || p.Lng != pos.Lng
            }).value()
            console.log("points selection: ", points)
            addNodeToSelect(points)
        }else{
            console.log("起点：",distributor.StartPos)
            console.log("终点：",distributor.DestPos)
            addNodeToSelect([distributor.StartPos, distributor.DestPos])
        }          
    }
    function addNodeToSelect(points){
        var options = $("#nodeSelect option").remove()
        points = _.remove(points, null)
        _.each(points, function(p){
            $("#nodeSelect").append(String.format('<option value="{0}">ID: {0}  ({1}, {2}) {3}</option>', p.ID, p.Lng, p.Lat, p.Address));
        })        
    }
    function isDistributorOnNode(dst){
        var crt = dst.CurrentPos
        if(crt != null && crt.ID > 0){
            return true
        }
        // var start = dst.StartPos
        // var dest = des.DestPos
        // if(crt == null) return false
        // if(crt.ID == start.ID){
        //     return true
        // }

        // if(crt.Lat == start.Lat && crt.Lng == start.Lng){
        //     return true
        // }
        return false
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
            conn.onmessage = onwsMessage
        } else {
            appendLog("浏览器不支持")
        }        
    }
    function send(msg){
        if (!conn) {
            return false;
        }
        console.log("send => ",msg)
        conn.send(JSON.stringify(msg))
    }
    //将地图设置为初始状态，目的是不突出任何信息
    function resetMap2Initial(){
        setMapMarker(116.404, 39.915, false)
    }
    function initMap(){
        // map.centerAndZoom(new BMap.Point(116.404, 39.915), 16);
        map.addControl(new BMap.ZoomControl());  //添加地图缩放控件
    }
    function setMapMarker(lng,lat, bAddMarker){
        map.removeOverlay(marker)
        if(lng > 0 || lat > 0){
            map.centerAndZoom(new BMap.Point(lng, lat), 16);
            // map.addControl(new BMap.ZoomControl());  //添加地图缩放控件
            if(bAddMarker == true){
                marker = new BMap.Marker(new BMap.Point(lng, lat));  //创建标注
                map.addOverlay(marker);                 // 将标注添加到地图中            
            }
        }
    }
    function hideOrderSelectButton(){
        // $("#btnAccept").hide()
        // $("#btnPass").hide()
        // $("#orderIDList").hide()
        // $("#btnStartDistribute").hide()
    }
    function hideAllControls(){
        // $("#btnPrepared").hide()
        // $("#btnAccept").hide()
        // $("#btnPass").hide()
        // $("#orderIDList").hide()
        // $("#btnStartDistribute").hide()
        // $("#allmap").hide()
    }
    function showOrderSelectButton(){
        // $("#btnAccept").show()
        // $("#btnPass").show()
        // $("#orderIDList").show()
        // $("#allmap").show()
    }
    function appendLog(message) {
        // if(output.children().length > 20){
        //     output.empty()
        // }
        output.prepend(message+"</br>");
    }    
    String.format = function(){
        if( arguments.length == 0){
            return null; 
        } 
        var str = arguments[0]; 
        for(var i=1;i<arguments.length;i++){
            var re = new RegExp('\\{' + (i-1) + '\\}','gm'); 
            str = str.replace(re, arguments[i]); 
        } 
        return str; 
    } 
   // function mapInit(){
   //     map = new BMap.Map("allmap");
   //     // var point = new BMap.Point(116.644691, 39.934758);//北京物资学院
   //     var point = new BMap.Point(116.331398,39.897445);//天安门
   //     map.centerAndZoom(point,12);
   //     map.addControl(new BMap.ZoomControl());  //添加地图缩放控件
   // }

   // function setMapMarker(lng,lat, bAddMarker){
   //     map.removeOverlay(marker)
   //     if(lng > 0 || lat > 0){
   //         map.centerAndZoom(new BMap.Point(lng, lat), 12);
   //         if(bAddMarker == true){
   //             marker = new BMap.Marker(new BMap.Point(lng, lat));  //创建标注
   //             map.addOverlay(marker);                 // 将标注添加到地图中            
   //         }
   //     }
   // } 
</script>

</html>
