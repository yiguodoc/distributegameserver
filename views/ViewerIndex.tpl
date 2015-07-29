<!DOCTYPE html>
<html lang="en">
<head>
<title>进度观看</title>
<script src="javascripts/jquery.js"></script>
<script src="javascripts/underscore.js"></script>
<script type="text/javascript" src="http://api.map.baidu.com/api?v=2.0&ak=kU4NWwyP5SwguC2W2WAfO1bO"></script>

<style type="text/css">
        body, html,#allmap {width: 100%;height: 100%;overflow: hidden;margin:0;font-family:"微软雅黑";}

</style>
</head>
<body>
    <div id="allmap"  style="height:70%;margin-top:20px;"></div>

    <div id="output"></div>

</body>
<script type="text/javascript">
    var conn;
    var output;
    var distributors = []
    var orders = []
    var map = null
    var markerMapOrderID = []
    $(function() {
        output = $("#output")
        initialMap()
        //加载数据完成后再建立连接

        $.get("/orders",function(data){//配送员数据
            console.log(data)
            orders = data
            console.log("共有 %d 个订单",_.size(orders))
            _.each(orders, function(order){
                var map = _.findWhere({OrderID: order.ID})
                if(map == null){
                    var pos = order.GeoSrc
                    addMapMarker(parseFloat(pos.Lng),parseFloat(pos.Lat),"grey", order.ID)
                }
            })
            $.get("/distributors",function(data){//配送员数据
                console.log(data)
                distributors = data
                console.log("共有 %d 个配送员",_.size(distributors))
                _.each(distributors, function(distributor){
                    // console.log(distributor.Name)
                    if(distributor.Online == true){
                        appendLog(distributor.Name +" 在线, 有 " + _.size(distributor.AcceptedOrders) + " 个订单")
                    }else{
                        appendLog(distributor.Name +" 离线, 有 " + _.size(distributor.AcceptedOrders) + " 个订单")
                    }
                    var acceptOrders = distributor.AcceptedOrders
                    _.each(acceptOrders, function(order){
                        var pos = order.GeoSrc
                        setMapMarker(distributor.Color,order.ID)
                    })
                })
            })            


            prepareConn()
        })

    })
    function initialMap(){
        map = new BMap.Map("allmap");
        var point = new BMap.Point(116.644691, 39.934758);//北京物资学院
        // var point = new BMap.Point(116.331398,39.897445);
        map.centerAndZoom(point,16);
        var top_left_navigation = new BMap.NavigationControl();  //左上角，添加默认缩放平移控件
        map.addControl(top_left_navigation);
        map.setMapStyle({style:"grayscale"}); 
    }
    function addMapMarker(lng,lat,icon, orderID){
        if(icon == null){
            icon = "grey"
        }
        if(lng > 0 || lat > 0){
            // map.centerAndZoom(new BMap.Point(lng, lat), 12);
            // map.setCenter(new BMap.Point(lng, lat));
            var imageUrl = "/images/marker/"+icon+".png"
            var myIcon = new BMap.Icon(imageUrl, new BMap.Size(20, 31), {anchor: new BMap.Size(10, 31)});
            var marker = new BMap.Marker(new BMap.Point(lng, lat), {icon: myIcon});  //创建标注
            map.addOverlay(marker);                 // 将标注添加到地图中
            markerMapOrderID.push({OrderID: orderID, Marker: marker})            
        }else{
            console.warn("坐标数据异常")
        }
    }
    function setMapMarker(icon, orderID){
        if(icon == null){
            icon = "grey"
        }
        var m = _.findWhere(markerMapOrderID, {OrderID: orderID})
        if(m == null){
            console.warn("没有找到订单对应的标记")
            return
        }
        var crtIcon = m.Marker.getIcon()
        if(crtIcon == null){
            console.error("系统异常")
            return
        }

        var imageUrl = "/images/marker/"+icon+".png"
        if(crtIcon.imageUrl == imageUrl){
            return
        }else{
            var myIcon = new BMap.Icon(imageUrl, new BMap.Size(36,94));
            m.Marker.setIcon(myIcon)
        }
    }
    function appendLog(message) {
        output.prepend(message+"</br>");
    }    
    function prepareConn(){
        if (window["WebSocket"]) {
            conn = new WebSocket("ws://{{.HOST}}/wsViewer?id={{.ID}}");
            conn.onclose = function(evt) {
                appendLog("连接关闭")
            }
            conn.onopen = function(evt){
                appendLog("连接成功")
            }
            conn.onmessage = function(evt) {
                // appendLog($("<div/>").text(evt.data))
                var msg = evt.data
                console.log(msg)
                msg = JSON.parse(msg)
                switch(msg.MessageType){
                    case 1://订单分发
                    
                    break
                    case 3://计时
                    appendLog("："+ msg.Data)
                    break
                    case 4://开始选择订单
                    break
                    case 5://消息广播
                    break
                    case 6://订单分配结果
                    var distributorID = msg.Data.DistributorID
                    var orderID = msg.Data.OrderID
                    var order = _.findWhere(orders, {ID: orderID})
                    var distributor = _.findWhere(distributors, {ID: distributorID})
                    if(order == null || distributor == null){
                        return
                    }else{
                        appendLog(distributor.Name + " 获取到订单 "+ order.ID +" 的配送权")
                        console.log("%s 获取到订单 %s 的配送权", distributor.Name, order.ID)
                        setMapMarker(distributor.Color, order.ID)
                    }
                    break
                    case 7://上线
                    console.log("上线")
                    console.log(msg.Data)
                    appendLog(msg.Data.Name + " 上线")
                    break
                    case 8://下线
                    console.log("下线")
                    console.log(msg.Data)
                    appendLog(msg.Data.Name + " 下线")
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
</script>

</html>
