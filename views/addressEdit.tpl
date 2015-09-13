<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="initial-scale=1.0, user-scalable=no" />
	<style type="text/css">
		body, html,#allmap {width: 100%;height: 100%;overflow: hidden;margin:0;font-family:"微软雅黑";}
	</style>
	<script src="javascripts/jquery.js"></script>
	<!-- // <script src="javascripts/underscore.js"></script> -->
	<script src="javascripts/lodash.js"></script>
	<script type="text/javascript" src="http://api.map.baidu.com/api?v=2.0&ak=kU4NWwyP5SwguC2W2WAfO1bO"></script>
	<script type="text/javascript" src="http://api.map.baidu.com/library/CurveLine/1.5/src/CurveLine.min.js"></script>
	<title>系统地图编辑器</title>
</head>
<body>
	<div style="margin-top:20px;margin-bottom:5px;">
        <input id="btnSelectMarker" type="button" value="选择点" onclick="switchControl(6)" style="margin-bottom: 10px;">
        <input id="btnAddMarker" type="button" value="添加点" onclick="switchControl(0)" style="margin-bottom: 10px;">
        <input id="btnRemoveMarker" type="button" value="删除点" onclick="switchControl(1)" style="margin-bottom: 10px;">
        <input id="btnRemoveMarker" type="button" value="设为路径节点" onclick="switchControl(10)" style="margin-bottom: 10px;">
        <input id="btnRemoveMarker" type="button" value="设为配送中心" onclick="switchControl(7)" style="margin-bottom: 10px;">
        <input id="btnRemoveMarker" type="button" value="设为出生点" onclick="switchControl(8)" style="margin-bottom: 10px;">
        <input id="btnRemoveMarker" type="button" value="设为非出生点" onclick="switchControl(9)" style="margin-bottom: 10px;">
        <input id="btnAddOrder" type="button" value="添加订单" onclick="switchControl(4)" style="margin-bottom: 10px;">
        <input id="btnRemoveOrder" type="button" value="移除订单" onclick="switchControl(5)" style="margin-bottom: 10px;">
        <input id="btnAddRoute" type="button" value="添加路径" onclick="switchControl(2)" style="margin-bottom: 10px;">
        <input id="btnRemoveRoute" type="button" value="移除路径" onclick="switchControl(3)" style="margin-bottom: 10px;">
        <input id="btnClearMapData" type="button" value="清除地图数据" onclick="clearMapData()" style="margin-bottom: 10px;">

	</div>
	<div id="allmap"  style="height:60%;"></div>
	<div id ="addressEditBox" style="margin-top: 20px;">
		<span>当前地址：</span><input id="address" type="text" value="" style="width:90%;"></br>
		<span>订单分值：</span><input id="orderScore" type="number" value="" style="width:90%;"></br>
		<span>地址坐标：</span><span id="lnglat"></span></br>
        <input id="btnSetAddress" type="button" value="保存" onclick="saveMarkerAddress()" style="">
	</div>
	<div style="margin-top:10px;margin-bottom:5px;">
        <input id="btnSaveData" type="button" value="保存地图设置" onclick="onSaveData()" style="margin-bottom: 10px;">
	</div>
</body>
</html>
<script type="text/javascript">

	var POSITION_TYPE_WAREHOUSE   = 0  //仓库
	var POSITION_TYPE_ORDER_ROUTE = 1  //路径节点
	var POSITION_TYPE_ORDER       = 2  //放置订单

    var map = null;
    //当前操作的选择
    //0 添加点  1 删除点 2 添加路径 3 移除路径 4 添加订单  5 移除订单 6 选择点 7 设为配送中心 8 设为出生点 9 设为非出生点  10 设为路径节点
    var optSelect = 6

    var markers = []
    var lines  = []
    var lineStartMarker = null
    var selectedMarker = null
    //各种类型的点的图标设计
	var bornPointIconDef = {isBornPoint: true, imageName: "aimRed.png", width: 100, height: 100, opt: {anchor: new BMap.Size(20, 20), imageSize: new BMap.Size(40,40)}}
	var iconKinds = [
		{pointType: POSITION_TYPE_WAREHOUSE, imageName: "warehouse.png", width: 64, height: 64, opt: {anchor: new BMap.Size(32, 48)}},
		{pointType: POSITION_TYPE_ORDER_ROUTE, imageName: "node.png", width: 52, height: 52, opt: {anchor: new BMap.Size(6, 6), imageSize: new BMap.Size(12,12)}},
		{pointType: POSITION_TYPE_ORDER, imageName: "bagageClosed.png", width: 29, height: 29, opt: {anchor: new BMap.Size(15, 15)}},
		bornPointIconDef
	]
	//出生点的图标
	var actionHandlerMap = [
		{optIndex: 0, receiverType: BMap.Map, actor: function(map, e){
			if(markerExist(e.point) == false){//不能重复
	 			var marker = addMapMarker(e.point, {Address: "", PointType: POSITION_TYPE_ORDER_ROUTE})
	 			markers.push(marker)	
			}
		}},
		{optIndex: 1, receiverType: BMap.Marker, actor: function(thisMarker){
			if(markerExist(thisMarker.getPosition()) == true){
				map.removeOverlay(thisMarker)
				removeMarker(thisMarker.getPosition())
			}	
		}},
		{optIndex: 2, receiverType: BMap.Marker, actor: function(thisMarker){
			if(lineStartMarker == null){//说明是起点
				lineStartMarker = thisMarker//那么将当前点击的marker的坐标作为线的起点
				// var label = new BMap.Label("我是文字标注哦",{offset:new BMap.Size(-5,-20)});
				thisMarker.setAnimation(BMAP_ANIMATION_BOUNCE); //跳动的动画
			}else{
				addLine(lineStartMarker.getPosition(),thisMarker.getPosition())
				lineStartMarker.setAnimation(null); //跳动的动画
				lineStartMarker = null
			}
		}},
		{optIndex: 3, receiverType: BMap.Polyline, actor: function(line){
			map.removeOverlay(line)
			removeLine(line.getPath())
		}},
		{optIndex: 4, receiverType: BMap.Marker, actor: function(thisMarker){
			thisMarker.pointType = POSITION_TYPE_ORDER
			resetMarkerIcon(thisMarker, {pointType: thisMarker.pointType})
		}},
		{optIndex: 5, receiverType: BMap.Marker, actor: function(thisMarker){
			thisMarker.pointType = POSITION_TYPE_ORDER_ROUTE
			resetMarkerIcon(thisMarker,  {pointType: thisMarker.pointType})
		}},
		{optIndex: 6, receiverType: BMap.Marker, actor: function(thisMarker, e){

			if(selectedMarker != null){
				selectedMarker.setAnimation(null)
			}
			selectedMarker = thisMarker
			resetAddressEditBoxStatus(true)
			onClickEditMarkAddress()
		}},
		{optIndex: 7, receiverType: BMap.Marker, actor: function(thisMarker){
			thisMarker.pointType = POSITION_TYPE_WAREHOUSE
			resetMarkerIcon(thisMarker, {pointType: thisMarker.pointType})
		}},
		{optIndex: 8, receiverType: BMap.Marker, actor: function(thisMarker){
			setMarkerBornPoint(thisMarker)
		}},
		{optIndex: 9, receiverType: BMap.Marker, actor: function(thisMarker){
			unsetMarkerBornPoint(thisMarker)
		}},
		{optIndex: 10, receiverType: BMap.Marker, actor: function(thisMarker){
			thisMarker.pointType = POSITION_TYPE_ORDER_ROUTE
			resetMarkerIcon(thisMarker,  {pointType: thisMarker.pointType})
		}}


	]
	$(function() {
		//初始化地图数据
		$.get("/mapData", function(data){
			console.log(data)
		    mapInit()
		    markers = _.reduce(data.Points, function(markerList, p){
    	    	var m =addMapMarker(new BMap.Point(p.Lng,p.Lat), {Address: p.Address, PointType: p.PointType})
    	    	// if(m != null){
    	    	// 	// m.address = p.Address
    	    	// 	// m.hasOrder = p.HasOrder
    	    	// 	// m.pointType = p.PointType
    	    	// 	// if(m.hasOrder == true){
    		    // 	// 	resetMarkerIcon(m, 3)
    	    	// 	// }else{
    		    // 	// 	resetMarkerIcon(m, m.pointType)
    	    	// 	// }
    	    	// 	resetMarkerIcon(m)
    	    	// }
    	    	markerList.push(m)
    	    	return markerList
		    }, [])
		    // _.forEach(data.Points, function(p){
		    // 	var m =addMapMarker(new BMap.Point(p.Lng,p.Lat))
		    // 	if(m != null){
		    // 		m.address = p.Address
		    // 		m.hasOrder = p.HasOrder
		    // 		m.pointType = p.PointType
		    // 		if(m.hasOrder == true){
			   //  		resetMarkerIcon(m, 3)
		    // 		}else{
			   //  		resetMarkerIcon(m, m.pointType)
		    // 		}
		    // 	}
		    // })
		    _.forEach(data.Lines, function(line, index){
		    	var l = addLine(new BMap.Point(line.Start.Lng,line.Start.Lat), new BMap.Point(line.End.Lng,line.End.Lat))
		    	if(l != null){
			    	console.log("line count: %d", index+1)
		    	}
		    })
		})
	})
	function mapInit(){
		map = new BMap.Map("allmap");
		var point = new BMap.Point(116.644691, 39.934758);//北京物资学院
		// var point = new BMap.Point(116.331398,39.897445);//天安门
		map.centerAndZoom(point,16);
		var top_left_navigation = new BMap.NavigationControl();  //左上角，添加默认缩放平移控件
		map.addControl(top_left_navigation);
 		map.addEventListener("click", optHandler);
	}
	// //点击地图时的处理函数，一般是添加新点
	// function mapClickHandler(e){
	// 	console.log("点击地图 %f,%f",e.point.lng, e.point.lat)
	// 	if(optSelect == 0){
	// 		if(markerExist(e.point) == false){//不能重复
	//  			var marker = addMapMarker(e.point, {Address: "", PointType: POSITION_TYPE_ORDER_ROUTE})
	//  			markers.push(marker)	
	// 		}
	// 	}		
	// }
	function unsetMarkerBornPoint(marker){
		if(marker.bornPoint != null){
			map.removeOverlay(marker.bornPoint)
			map.bornPoint = null
		}
	}
	function setMarkerBornPoint(marker){
		var pos = marker.getPosition()
		marker.bornPoint = new BMap.Marker(pos);
		resetMarkerIcon(marker.bornPoint, {isBornPoint: true})
		map.addOverlay(marker.bornPoint)
		marker.setTop(true)
	}
	function resetMarkerIcon(marker, selectOpt){
		var iconDef = _.findWhere(iconKinds, selectOpt)
		// var iconDef = _.findWhere(iconKinds, {pointType: marker.pointType})
		if(iconDef != null){
			var imageUrl = "/images/marker/"+ iconDef.imageName
			var myIcon = new BMap.Icon(imageUrl, new BMap.Size(iconDef.width, iconDef.height), iconDef.opt);
			marker.setIcon(myIcon)		
		}else{
			console.error("cannot find icon type",selectOpt)
		}

		// var myIcon
		// switch(pointType){
		// 	case 0:
		//     	console.log("设置marker 为 仓库")
		// 		myIcon = createIcon("warehouse")
		// 	break
		// 	case 1:
		//     	console.log("设置marker 为 路径节点")
		// 		myIcon = createIcon("keyPoint")
		// 	break
		// 	case 2:
		//     	console.log("设置marker 为 途经点")
		// 		myIcon = createIcon("passPoint")
		// 	break
		// 	case 3:
		//     	console.log("设置marker 为 带有订单的途经点")
		// 		myIcon = createIcon("order")				
		// 	break
		// }
		// marker.setIcon(myIcon)		
	}
	// //右键菜单的处理函数
	// //每次根据点类型重置右键菜单
	// function markerMenuHandler(markerMenu, pointType, e, ee){
	// 	var marker = this
	// 	marker.pointType = pointType
	// 	resetMarkerIcon(marker)
	// 	marker.removeContextMenu(markerMenu)
 //    	marker.addContextMenu(createContextMenu(marker, pointType));		
	// }
	// function createContextMenu(marker, pointType){
	// 	var markerMenu = new BMap.ContextMenu();
 //    	switch(pointType){
 //    		case POSITION_TYPE_WAREHOUSE:
 //    		markerMenu.addItem(new BMap.MenuItem('<div style="font-size:16px;padding-top: 5px;">设为途经点</div>',markerMenuHandler.bind(marker, markerMenu, 2)));
 //    		markerMenu.addItem(new BMap.MenuItem('<div style="font-size:16px;padding-bottom: 5px;">设为路径节点</div>',markerMenuHandler.bind(marker, markerMenu, 1)));
 //    		break
 //    		case POSITION_TYPE_ORDER_ROUTE:
 //    		markerMenu.addItem(new BMap.MenuItem('<div style="font-size:16px;padding-top: 5px;">设为途经点</div>',markerMenuHandler.bind(marker, markerMenu, 2)));
 //    		markerMenu.addItem(new BMap.MenuItem('<div style="font-size:16px;padding-bottom: 5px;">设为仓库</div>',markerMenuHandler.bind(marker, markerMenu, 0)));
 //    		break
 //    		case POSITION_TYPE_ORDER:
 //    		markerMenu.addItem(new BMap.MenuItem('<div style="font-size:16px;padding-top: 5px;">设为仓库</div>',markerMenuHandler.bind(marker, markerMenu, 0)));
 //    		markerMenu.addItem(new BMap.MenuItem('<div style="font-size:16px;padding-bottom: 5px;">设为路径节点</div>',markerMenuHandler.bind(marker, markerMenu, 1)));
 //    		break
 //    	}
 //    	return markerMenu
	// }

	function optHandler(e){
    	// console.log(e)
    	// console.log(this)
    	console.log("当前操作 %d 点击", optSelect)
    	if(e && e.stopPropagation){
	    	e.stopPropagation()
    	}else{
	    	window.event.cancelBubble = true;//阻止事件冒泡，防止在点击marker时同时添加一个marker
    	}

    	var opt = _.findWhere(actionHandlerMap, {optIndex: optSelect})
    	if(opt != null){
    		if(this instanceof opt.receiverType && opt.actor != null){
	    		opt.actor(this,e)
    		}
    	}else{
    		console.error("系统异常操作")
    	}		
	}

	//添加一个marker到地图
	//bmapPoint：地图的坐标
	//opt:marker：附带的属性，根据属性确定marker的图标信息
	function addMapMarker(bmapPoint, opt){
		if(markerExist(bmapPoint) == false){
			// var myIcon = createIcon("passPoint")
		    var marker = new BMap.Marker(bmapPoint);  //创建标注
		    if(opt != null){
		    	if(_.has(opt, "PointType")){
				    marker.pointType = opt.PointType //默认为途经点
		    	}
		    	
		    	if(_.has(opt, "Address")){
				    marker.address = opt.Address
		    	}		    	
		    }
    		resetMarkerIcon(marker, {pointType: marker.pointType})
		    map.addOverlay(marker);                 // 将标注添加到地图中

		    marker.addEventListener("click", optHandler)
	    	//创建右键菜单
	    	// marker.addContextMenu(createContextMenu(marker, 2));
		    // marker.addEventListener("click", function(e){
		    	// var thisMarker = this
		    	// // console.log(e)
		    	// // console.log(this)
		    	// console.log("当前操作 %d 点击 marker", optSelect)
		    	// if(e && e.stopPropagation){
			    // 	e.stopPropagation()
		    	// }else{
			    // 	window.event.cancelBubble = true;//阻止事件冒泡，防止在点击marker时同时添加一个marker
		    	// }
		    	// optHandler(optSelect, thisMarker)
		    	// switch(optSelect){
		    	// 	case 1://删除点
				   //  	if(markerExist(thisMarker.getPosition()) == true){
				   //  		map.removeOverlay(thisMarker)
				   //  		removeMarker(thisMarker.getPosition())
				   //  	}	
		    	// 	break
		    	// 	case 2://添加路径
		    	// 		if(lineStartMarker == null){//说明是起点
		    	// 			lineStartMarker = thisMarker//那么将当前点击的marker的坐标作为线的起点
		    	// 			// var label = new BMap.Label("我是文字标注哦",{offset:new BMap.Size(-5,-20)});
		    	// 			thisMarker.setAnimation(BMAP_ANIMATION_BOUNCE); //跳动的动画
		    	// 		}else{
		    	// 			addLine(lineStartMarker.getPosition(),thisMarker.getPosition())
		    	// 			lineStartMarker.setAnimation(null); //跳动的动画
	    		// 			lineStartMarker = null
		    	// 		}
		    	// 	break
		    	// 	case 4://添加订单
		    	// 		if(thisMarker.pointType != 1){
		    	// 			console.info("不是路径节点，无法添加订单")
		    	// 			alert("订单只能添加到路径节点上，请先将该点转换为路径节点")
		    	// 		}else{
		    	// 			thisMarker.hasOrder = true
		    	// 			resetMarkerIcon(thisMarker, 3)
		    	// 		}
		    	// 	break
		    	// 	case 5://移除订单
			    // 		if(thisMarker.pointType != 1){
			    // 			console.warn("在非路径节点上移除订单，有异常")
			    // 		}else{
			    // 			thisMarker.hasOrder = false
		    	// 			resetMarkerIcon(thisMarker, 1)
			    // 		}
		    	// 	break
		    	// 	case 6://选择点
		    	// 		if(selectedMarker != null){
		    	// 			selectedMarker.setAnimation(null)
		    	// 		}
		    	// 		selectedMarker = thisMarker
		    	// 		resetAddressEditBoxStatus(true)
		    	// 		onClickEditMarkAddress()
		    	// 	break
		    	// }
		    // })
		    // markers.push(marker)	
		    return marker		
		}
		return null
	}

	function addLine(startPoint, destPoint){
		if(destPoint.equals(startPoint) == true){//同一个点不能画线
			return null
		}
		var points = [startPoint, destPoint];
		if(lineExist(points) == true){
			return null
		}
		var line = new BMap.Polyline(points, {strokeColor:"blue", strokeWeight:5, strokeOpacity:0.5}); //创建弧线对象
		var distance = map.getDistance(startPoint, destPoint)//米
		line.Distance = distance
		map.addOverlay(line); //添加到地图中
		lines.push(line)
		line.addEventListener("click", optHandler)

		// line.addEventListener("click", function(e){
		// 	if(optSelect == 3){
		// 		map.removeOverlay(this)
		// 		removeLine(this.getPath())
		// 	}
		// })
		return line
	}
	function switchControl(opt){
		optSelect = opt
    	console.log("当前操作 %d", optSelect)
    	clearAddRouteStates()
	}
	function clearMapData(){
		var r = confirm("将清除所有的点和路径，不可恢复，确定吗？")
		if(r){
			_.forEach(markers, function(m){
				map.removeOverlay(m)
			})
			markers = []
			_.forEach(lines, function(line){
				map.removeOverlay(line)
			})
			lines=[]
		}
	}
	function onSaveData(){
		console.log("保存地图数据")
		//保存两类数据，点和线

		// console.log(linesData)
		var pointsData = _.map(markers, function(marker,index){
			var p = marker.getPosition()
			return {ID: index+1, Lat: p.lat, Lng: p.lng, PointType: marker.pointType, Address: marker.address, IsBornPoint: marker.bornPoint != null, Score: marker.score}
		})

		var linesData = _.map(lines, function(line){
			var points = line.getPath()
			var start = _.find(pointsData, function(pnt){
				return points[0].lat == pnt.Lat && pnt.Lng == points[0].lng
			})
			var stop = _.find(pointsData, function(pnt){
				return points[1].lat == pnt.Lat && pnt.Lng == points[1].lng
			})
			return {Start: start, End: stop, Distance: line.Distance}
			// return {Start: {Lat: points[0].lat, Lng: points[0].lng}, End: {Lat: points[1].lat, Lng: points[1].lng}, Distance: line.Distance}
		})		
		
		var uploadMapData = {Points: pointsData, Lines: linesData}
		console.log(uploadMapData)
		console.log(JSON.stringify(uploadMapData))
		$.post("/uploadMapData", {data: JSON.stringify(uploadMapData)}, function(data, status){
			console.log(data)
			if(data.Code == 0){
				alert("保存成功")
			}else{
				alert(data.Message)
			}
			// console.log(status)
		})
	}
	function onClickEditMarkAddress(){
		if(selectedMarker == null) return

		selectedMarker.setAnimation(BMAP_ANIMATION_BOUNCE); //跳动的动画
		var inputAddress = $("#address")
		var inputScore = $("#orderScore")
		var lnglat = $("#lnglat")

			// lnglat.val("")
		var p = selectedMarker.getPosition()
		lnglat.text("("+p.lng+", "+p.lat+")")
		if(selectedMarker.address == null){
			inputAddress.val("")
		}else{
			inputAddress.val(selectedMarker.address)
		}
		if(selectedMarker.score == null){
			inputScore.val(0)
		}else{
			inputScore.val(selectedMarker.score)
		}
	}
	function saveMarkerAddress(){
		if(selectedMarker == null) return

		var val = $("#address").val()
		var score = $("#orderScore").val()
		selectedMarker.address = val
		selectedMarker.score = score
		alert("设置该点的地址为 "+ val + "  分值为 " + score)
	}
	//设置地址编辑区域的显示或者隐藏
	function resetAddressEditBoxStatus(status){
		if(status){
			$("#addressEditBox").show()
		}else{
			$("#addressEditBox").hide()
		}
	}
	//清除因为添加连线而设置的状态，例如动画
	function clearAddRouteStates(){
		if(lineStartMarker != null){
			lineStartMarker.setAnimation(null); //跳动的动画
		}

		if(selectedMarker !=null){
			selectedMarker.setAnimation(null)
		}
		resetAddressEditBoxStatus(false)
	}
	function lineExist(points){
		line = _.find(lines, function(l){//排除包含了参数中所有点的线
			return _.every(points, function(p){
				return _.contains(l.getPath(), p)
			})
		})
		if(line == null){
			console.info("添加新的路径")
			return false
		}else{
			console.log("路径已经存在")
			return true
		}
	}
	function removeLine(points){
		lines = _.reject(lines, function(l){//排除包含了参数中所有点的线
			return _.every(points, function(p){
				return _.contains(l.getPath(), p)
			})
		})
		console.log("剩余 %d 条路径", _.size(lines))
	}
	function removeMarker(bmapPoint){
		markers = _.reject(markers, function(m){
			return m.getPosition().equals(bmapPoint)
		})
		console.log("剩余 %d 个 Marker", _.size(markers))
	}
	function markerExist(bmapPoint){
		var marker = _.find(markers, function(m){
			// console.log(m.getPosition())
			// console.log(bmapPoint)
			return m.getPosition().equals(bmapPoint)
		})
		return marker!=null
	}
	// // 创建地址解析器实例
	// var myGeo = new BMap.Geocoder();

	// map.addEventListener("click",function(e){
	// 		// alert(e.point.lng + "," + e.point.lat);
	// 		console.log("%f,%f",e.point.lng, e.point.lat)
	// 		myGeo.getLocation(e.point, function(rs){
	// 			var addComp = rs.addressComponents;
	// 			console.log(rs)
	// 			// appendLog(rs.address  +"  " + e.point.lng + ", "+ e.point.lat)
	// 			// alert(addComp.province + ", " + addComp.city + ", " + addComp.district + ", " + addComp.street + ", " + addComp.streetNumber);
	// 		}); 

			
	// });

	// // 将地址解析结果显示在地图上,并调整地图视野
	// myGeo.getPoint("北京市通州区北京物资学院", function(point){
	// 	if (point) {
	// 		console.log(point)
	// 		map.centerAndZoom(point, 16);
	// 		map.addOverlay(new BMap.Marker(point));
	// 	}else{
	// 		alert("您选择地址没有解析到结果!");
	// 	}
	// }, "北京市");

	function appendLog(message) {
	    // output.prepend(message+"</br>");
	    console.log(message)
	}    
</script>