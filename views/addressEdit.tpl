<!DOCTYPE html>
<html style="height:100%;">

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no" />
    <link rel="icon" type="image/png" href="/images/logo.png">
    <link href="http://g.alicdn.com/sj/dpl/1.5.1/css/sui.min.css" rel="stylesheet">
    <!-- <link href="stylesheets/docs.css" rel="stylesheet"> -->
    <style type="text/css">
    /*body, html,#allmap {width: 100%;height: 100%;overflow: hidden;margin:0;font-family:"微软雅黑";}*/
    </style>
    <script src="javascripts/jquery.js"></script>
    <!-- // <script type="text/javascript" src="http://g.alicdn.com/sj/lib/jquery/dist/jquery.min.js"></script> -->
    <script type="text/javascript" src="http://g.alicdn.com/sj/dpl/1.5.1/js/sui.min.js"></script>
    <!-- // <script src="javascripts/application.js"></script> -->
    <!-- // <script src="javascripts/underscore.js"></script> -->
    <script src="javascripts/lodash.js"></script>
    <script type="text/javascript" src="http://api.map.baidu.com/api?v=2.0&ak=kU4NWwyP5SwguC2W2WAfO1bO"></script>
    <script type="text/javascript" src="http://api.map.baidu.com/library/CurveLine/1.5/src/CurveLine.min.js"></script>
    <title>系统地图编辑器</title>
</head>

<body style="height:100%;">
    <div class='container'>
        <div class="sui-navbar navbar-inverse">
            <div class="navbar-inner" style="height:60px;">
                <div class="sui-container" style="margin-top: 10px; font-size: 14px;margin-left: 32px;"><a href="#" class="sui-brand">配送大师</a>
                    <ul class="sui-nav" style="margin-left: 28px;">
                        <li class="active"><a href="#">首页</a></li>
                        <!-- <li><a href="#">组件</a></li> -->
                        <li class="sui-dropdown"><a href="javascript:void(0);" data-toggle="dropdown" class="dropdown-toggle">其他 <i class="caret"></i></a>
                            <ul role="menu" class="sui-dropdown-menu">
                                <li role="presentation"><a role="menuitem" tabindex="-1" href="#">关于</a></li>
                                <li role="presentation"><a role="menuitem" tabindex="-1" href="#">项目组成员</a></li>
                                <li role="presentation"><a role="menuitem" tabindex="-1" href="#">版权</a></li>
                            </ul>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
    <div class="sui-layout" style="height:85%;margin-top:0px;">
        <div class="sidebar" style="position: static;">
            <div style="text-align:center;border-top: 1px solid rgba(100,100,100,0.3);">
                <div style="  margin-bottom: 20px; margin-top: 10px; font-size: 16px; border-bottom: solid 1px rgba(100,100,100,0.2); padding-bottom: 10px; color: rgba(100,100,100,0.7);">功能区</div>
                <div style="margin-bottom: 0px;">
                    <div style="font-size: 14px; margin-bottom: 5px;">选择地图</div>
                    <select id="selectMap" style="width: 70%; height: 22px;">
                        <!-- <option value="saab">Saab</option> -->
                    </select>
                    <a href="javascript:void(0);" onclick="loadSelectedMap()" class="sui-btn  btn-info" style="width:70%;margin-bottom:0px;margin-top: 10px;">下载</a></br>
                </div>
                <div style="border-bottom: solid 1px rgba(100,100,100,0.2);margin-bottom: 15px; margin-top: 15px;"></div>
                <a href="javascript:void(0);" onclick="switchControl(6)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">选择点</a></br>
                <a href="javascript:void(0);" onclick="switchControl(0)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">添加点</a></br>
                <a href="javascript:void(0);" onclick="switchControl(1)" class="sui-btn  btn-danger" style="width:70%;margin-bottom:10px;">删除点</a></br>
                
                 <div style="width:70%;border-bottom: solid 1px rgba(100,100,100,0.3);margin-bottom: 10px; margin-top: 5px;margin-left:15%;"></div>
                
                <a href="javascript:void(0);" onclick="switchControl(2)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">添加路径</a></br>
                <a href="javascript:void(0);" onclick="switchControl(3)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">移除路径</a></br>

                 <div style="width:70%;border-bottom: solid 1px rgba(100,100,100,0.3);margin-bottom: 10px; margin-top: 5px;margin-left:15%;"></div>

                <a href="javascript:void(0);" onclick="switchControl(10)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">设为路径节点</a></br>
                <a href="javascript:void(0);" onclick="switchControl(7)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">设为配送中心</a></br>
                <a href="javascript:void(0);" onclick="switchControl(8)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">设为出生点</a></br>
                <a href="javascript:void(0);" onclick="switchControl(9)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">设为非出生点</a></br>
                
                 <div style="width:70%;border-bottom: solid 1px rgba(100,100,100,0.3);margin-bottom: 10px; margin-top: 5px;margin-left:15%;"></div>

                <a href="javascript:void(0);" onclick="switchControl(4)" class="sui-btn  btn-info" style="width:70%;margin-bottom:10px;">添加订单</a></br>
                <a href="javascript:void(0);" onclick="switchControl(5)" class="sui-btn  btn-danger" style="width:70%;margin-bottom:10px;">移除订单</a></br>
                
                 <div style="width:70%;border-bottom: solid 1px rgba(100,100,100,0.3);margin-bottom: 10px; margin-top: 5px;margin-left:15%;"></div>

                <a href="javascript:void(0);" onclick="clearMapData()" class="sui-btn btn-danger" style="width:70%;margin-bottom:10px;">清除地图数据</a></br>

                 <div style="width:70%;border-bottom: solid 1px rgba(100,100,100,0.3);margin-bottom: 10px; margin-top: 5px;margin-left:15%;"></div>

                <div style="font-size: 14px; margin-bottom: 5px;">游戏时长（分钟）</div>
                <input id="gameTimeLength" type="number" value="" style="width:70%;text-align:center;">


                <a href="javascript:void(0);" onclick="onSaveData()" class="sui-btn btn-xlarge btn-success" style="width:70%;margin-top:20px;margin-bottom:10px;">保存地图</a></br>
            </div>
        </div>
        <div class="content" style="height:100%;margin-left: 195px; margin-right: 5px;border-left: 3px solid rgba(100,100,100,0.3); padding-left: 2px;">
            <div id="allmap" style="height:80%;border-top: 1px solid rgba(100,100,100,0.3); border-bottom: 1px solid rgba(100,100,100,0.3);"></div>
            <div id="addressEditBox" style="margin-top: 10px;">
                <div style="margin-bottom:10px;">
                    <span>当前地址：</span>
                    <input id="address" type="text" value="" style="width:98%;">
                </div>
                <div style="margin-bottom:10px;">
                    <span>订单分值：</span>
                    <input id="orderScore" type="number" value="" style="width:98%;">
                    </br>
                </div>
                <!-- <span>地址坐标：</span><span id="lnglat"></span></br> -->
                <!-- <input id="btnSetAddress" type="button" value="保存" onclick="saveMarkerAddress()" style=""> -->
                <a href="javascript:void(0);" onclick="saveMarkerAddress()" class="sui-btn  btn-info" style="margin-left:5px; width:100px;">保存</a>
            </div>
            <!-- 			<div style="margin-top:10px;margin-bottom:5px;">
			    <input id="btnSaveData" type="button" value="保存地图设置" onclick="onSaveData()" style="margin-bottom: 10px;">
			</div>
 -->
        </div>
    </div>
    <!-- <div style="width: 100%; text-align: center; font-size: 13px; padding-top: 10px; border-top: 1px solid rgba(100,100,100,0.3); color: rgba(100,100,100,0.8); margin-top: 5px;">配送大师团队技术支持</div> -->
    <script type="text/javascript">
    var POSITION_TYPE_WAREHOUSE = 0 //仓库
    var POSITION_TYPE_ORDER_ROUTE = 1 //路径节点
    var POSITION_TYPE_ORDER = 2 //放置订单

    var map = null;
    //当前操作的选择
    //0 添加点  1 删除点 2 添加路径 3 移除路径 4 添加订单  5 移除订单 6 选择点 7 设为配送中心 8 设为出生点 9 设为非出生点  10 设为路径节点
    var optSelect = 6

    var markers = []
    var lines = []
    var lineStartMarker = null
    var selectedMarker = null
    var currentMapID = null//当前地图名称
        //各种类型的点的图标设计
    var bornPointIconDef = {
        isBornPoint: true,
        imageName: "aimRed.png",
        width: 100,
        height: 100,
        opt: {
            anchor: new BMap.Size(20, 20),
            imageSize: new BMap.Size(40, 40)
        }
    }
    var iconKinds = [{
                pointType: POSITION_TYPE_WAREHOUSE,
                imageName: "warehouse.png",
                width: 64,
                height: 64,
                opt: {
                    anchor: new BMap.Size(24, 24),
                    imageSize: new BMap.Size(48, 48)
                }
            }, {
                pointType: POSITION_TYPE_ORDER_ROUTE,
                imageName: "node.png",
                width: 52,
                height: 52,
                opt: {
                    anchor: new BMap.Size(6, 6),
                    imageSize: new BMap.Size(12, 12)
                }
            }, {
                pointType: POSITION_TYPE_ORDER,
                imageName: "bagageClosed.png",
                width: 29,
                height: 29,
                opt: {
                    anchor: new BMap.Size(15, 15)
                }
            },
            bornPointIconDef
        ]
        //出生点的图标
    var actionHandlerMap = [
    	{
	        optIndex: 0,
	        receiverType: BMap.Map,
	        actor: function(map, e) {
	            if (markerExist(e.point) == false) { //不能重复
	                var marker = addMapMarker(e.point, {
	                    Address: "",
	                    PointType: POSITION_TYPE_ORDER_ROUTE
	                })
	                markers.push(marker)
	            }
	        }
	    }, {
	        optIndex: 1,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            if (markerExist(thisMarker.getPosition()) == true) {
	                map.removeOverlay(thisMarker)
	                removeMarker(thisMarker.getPosition())
	            }
	        }
	    }, {
	        optIndex: 2,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            if (lineStartMarker == null) { //说明是起点
	                lineStartMarker = thisMarker //那么将当前点击的marker的坐标作为线的起点
	                    // var label = new BMap.Label("我是文字标注哦",{offset:new BMap.Size(-5,-20)});
	                thisMarker.setAnimation(BMAP_ANIMATION_BOUNCE); //跳动的动画
	            } else {
	                addLine(lineStartMarker.getPosition(), thisMarker.getPosition())
	                lineStartMarker.setAnimation(null); //跳动的动画
	                lineStartMarker = null
	            }
	        }
	    }, {
	        optIndex: 3,
	        receiverType: BMap.Polyline,
	        actor: function(line) {
	            map.removeOverlay(line)
	            removeLine(line.getPath())
	        }
	    }, {
	        optIndex: 4,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            thisMarker.pointType = POSITION_TYPE_ORDER
	            resetMarkerIcon(thisMarker, {
	                pointType: thisMarker.pointType
	            })
	        }
	    }, {
	        optIndex: 5,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            thisMarker.pointType = POSITION_TYPE_ORDER_ROUTE
	            resetMarkerIcon(thisMarker, {
	                pointType: thisMarker.pointType
	            })
	        }
	    }, {
	        optIndex: 6,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker, e) {

	            if (selectedMarker != null) {
	                selectedMarker.setAnimation(null)
	            }
	            selectedMarker = thisMarker
	            resetAddressEditBoxStatus(true)
	            onClickEditMarkAddress()
	        }
	    }, {
	        optIndex: 7,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            thisMarker.pointType = POSITION_TYPE_WAREHOUSE
	            resetMarkerIcon(thisMarker, {
	                pointType: thisMarker.pointType
	            })
	        }
	    }, {
	        optIndex: 8,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            setMarkerBornPoint(thisMarker)
	        }
	    }, {
	        optIndex: 9,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            unsetMarkerBornPoint(thisMarker)
	        }
	    }, {
	        optIndex: 10,
	        receiverType: BMap.Marker,
	        actor: function(thisMarker) {
	            thisMarker.pointType = POSITION_TYPE_ORDER_ROUTE
	            resetMarkerIcon(thisMarker, {
	                pointType: thisMarker.pointType
	            })
	        }
	    }]
    $(function() {
        mapInit()
        $.get("/mapNameList", function(data) {
                console.log("map list => ", data.Data)
                var mapNameList = data.Data
                var $selectMapList = $("#selectMap")
                _.each(mapNameList, function(name) {
                    $selectMapList.append('<option value="' + name + '">' + name + '</option>')
                })
            })
    })
	
    function mapInit() {
        map = new BMap.Map("allmap");
        var point = new BMap.Point(116.644691, 39.934758); //北京物资学院
        // var point = new BMap.Point(116.331398,39.897445);//天安门
        map.centerAndZoom(point, 16);
        var top_right_navigation = new BMap.NavigationControl({
            anchor: BMAP_ANCHOR_TOP_RIGHT
        }); //左上角，添加默认缩放平移控件
        map.addControl(top_right_navigation);
        map.addEventListener("click", optHandler);
    }
    function loadMapData(mapID){
    	//初始化地图数据
    	$.get("/mapData?id="+mapID, function(data) {
	        console.log(data.Data)
	        var mapInfo = data.Data
	        markers = _.reduce(mapInfo.Points, function(markerList, p) {
	            var m = addMapMarker(new BMap.Point(p.Lng, p.Lat), {
	                Address: p.Address,
	                PointType: p.PointType,
	                Score: p.Score
	            })
	            markerList.push(m)
	            if (p.IsBornPoint) {
	                setMarkerBornPoint(m)
	            }
	            return markerList
	        }, [])

    	    _.forEach(mapInfo.Lines, function(line, index) {
    	        var l = addLine(new BMap.Point(line.Start.Lng, line.Start.Lat), new BMap.Point(line.End.Lng, line.End.Lat))
    	        if (l != null) {
    	            console.log("line count: %d", index + 1)
    	        }
    	    })
            $("#gameTimeLength").val(mapInfo.TimeLength)
    	    currentMapID = mapID
    	})
    }

    function unsetMarkerBornPoint(marker) {
        if (marker.bornPoint != null) {
            map.removeOverlay(marker.bornPoint)
            marker.bornPoint = null
        }
    }

    function setMarkerBornPoint(marker) {
        var pos = marker.getPosition()
        marker.bornPoint = new BMap.Marker(pos);
        resetMarkerIcon(marker.bornPoint, {
            isBornPoint: true
        })
        map.addOverlay(marker.bornPoint)
        marker.setTop(true)
    }

    function resetMarkerIcon(marker, selectOpt) {
        var iconDef = _.findWhere(iconKinds, selectOpt)
            // var iconDef = _.findWhere(iconKinds, {pointType: marker.pointType})
        if (iconDef != null) {
            var imageUrl = "/images/marker/" + iconDef.imageName
            var myIcon = new BMap.Icon(imageUrl, new BMap.Size(iconDef.width, iconDef.height), iconDef.opt);
            marker.setIcon(myIcon)
        } else {
            console.error("cannot find icon type", selectOpt)
        }
    }

    function optHandler(e) {
        console.log("当前操作 %d 点击", optSelect)
        if (e && e.stopPropagation) {
            e.stopPropagation()
        } else {
            window.event.cancelBubble = true; //阻止事件冒泡，防止在点击marker时同时添加一个marker
        }

        var opt = _.findWhere(actionHandlerMap, {
            optIndex: optSelect
        })
        if (opt != null) {
            if (this instanceof opt.receiverType && opt.actor != null) {
                opt.actor(this, e)
            }
        } else {
            console.error("系统异常操作")
        }
    }

    //添加一个marker到地图
    //bmapPoint：地图的坐标
    //opt:marker：附带的属性，根据属性确定marker的图标信息
    function addMapMarker(bmapPoint, opt) {
        if (markerExist(bmapPoint) == false) {
            // var myIcon = createIcon("passPoint")
            var marker = new BMap.Marker(bmapPoint); //创建标注
            if (opt != null) {
                if (_.has(opt, "PointType")) {
                    marker.pointType = opt.PointType //默认为途经点
                }

                if (_.has(opt, "Address")) {
                    marker.address = opt.Address
                }

                if (_.has(opt, "Score")) {
                    marker.score = opt.Score
                }
            }
            resetMarkerIcon(marker, {
                pointType: marker.pointType
            })
            map.addOverlay(marker); // 将标注添加到地图中

            marker.addEventListener("click", optHandler)

            return marker
        }
        return null
    }

    function addLine(startPoint, destPoint) {
        if (destPoint.equals(startPoint) == true) { //同一个点不能画线
            return null
        }
        var points = [startPoint, destPoint];
        if (lineExist(points) == true) {
            return null
        }
        var line = new BMap.Polyline(points, {
            strokeColor: "blue",
            strokeWeight: 5,
            strokeOpacity: 0.5
        }); //创建弧线对象
        var distance = map.getDistance(startPoint, destPoint) //米
        line.Distance = distance
        map.addOverlay(line); //添加到地图中
        lines.push(line)
        line.addEventListener("click", optHandler)
        return line
    }

    function switchControl(opt) {
        optSelect = opt
        console.log("当前操作 %d", optSelect)
        clearAddRouteStates()
    }

    function loadSelectedMap() {
        var $options = $("#selectMap option:selected")
            // console.log($options)
        if ($options.length <= 0) {
            alert("需要先选择地图")
            return
        }
        if(currentMapID != null){
        	clearMapData()
        }
        var selectedOption = $options[0]
        console.log("download map => " + selectedOption.value)
        loadMapData(selectedOption.value)
    }

    function clearMapData() {
    	_.forEach(markers, function(m) {
    	    map.removeOverlay(m)
    	})
    	markers = []
    	_.forEach(lines, function(line) {
    	    map.removeOverlay(line)
    	})
    	lines = []
        // var r = confirm("将清除所有的点和路径，不可恢复，确定吗？")
        // if (r) {
        //     _.forEach(markers, function(m) {
        //         map.removeOverlay(m)
        //     })
        //     markers = []
        //     _.forEach(lines, function(line) {
        //         map.removeOverlay(line)
        //     })
        //     lines = []
        // }
    }

    function onSaveData() {
    	if(currentMapID == null || currentMapID.length <= 0){
    		console.info("map ID wrong")
    		return
    	}
        console.log("保存地图数据")
            //保存两类数据，点和线
        var i = 0
            // console.log(linesData)
        var pointsData = _.map(markers, function(marker, index) {
            if (marker.bornPoint != null) {
                i++
                console.info("发现 %d 个出生点", i)
            }
            var p = marker.getPosition()
            return {
                ID: index + 1,
                Lat: p.lat,
                Lng: p.lng,
                PointType: marker.pointType,
                Address: marker.address,
                IsBornPoint: marker.bornPoint != null,
                Score: marker.score
            }
        })

        var linesData = _.map(lines, function(line) {
            var points = line.getPath()
            var start = _.find(pointsData, function(pnt) {
                return points[0].lat == pnt.Lat && pnt.Lng == points[0].lng
            })
            var stop = _.find(pointsData, function(pnt) {
                return points[1].lat == pnt.Lat && pnt.Lng == points[1].lng
            })
            return {
                Start: start,
                End: stop,
                Distance: line.Distance
            }
            // return {Start: {Lat: points[0].lat, Lng: points[0].lng}, End: {Lat: points[1].lat, Lng: points[1].lng}, Distance: line.Distance}
        })

        var uploadMapData = {
            Points: pointsData,
            Lines: linesData,
            TimeLength: parseInt($("#gameTimeLength").val())
        }
        console.log(uploadMapData)
        console.log(JSON.stringify(uploadMapData))
        $.post("/uploadMapData?id="+currentMapID, {
            data: JSON.stringify(uploadMapData)
        }, function(data, status) {
            console.log(data)
            if (data.Code == 0) {
                alert("保存成功")
            } else {
                alert(data.Message)
            }
            // console.log(status)
        })
    }

    function onClickEditMarkAddress() {
        if (selectedMarker == null) return

        selectedMarker.setAnimation(BMAP_ANIMATION_BOUNCE); //跳动的动画
        var inputAddress = $("#address")
        var inputScore = $("#orderScore")
        var lnglat = $("#lnglat")

        // lnglat.val("")
        var p = selectedMarker.getPosition()
        lnglat.text("(" + p.lng + ", " + p.lat + ")")
        if (selectedMarker.address == null) {
            inputAddress.val("")
        } else {
            inputAddress.val(selectedMarker.address)
        }
        if (selectedMarker.score == null) {
            inputScore.val(0)
        } else {
            inputScore.val(selectedMarker.score)
        }
    }

    function saveMarkerAddress() {
        if (selectedMarker == null) return

        var val = $("#address").val()
        var score = $("#orderScore").val()
        selectedMarker.address = val
        selectedMarker.score = parseInt(score)
        alert("设置该点的地址为 " + val + "  分值为 " + score)
    }
    //设置地址编辑区域的显示或者隐藏
    function resetAddressEditBoxStatus(status) {
        if (status) {
            $("#addressEditBox").show()
        } else {
            $("#addressEditBox").hide()
        }
    }
    //清除因为添加连线而设置的状态，例如动画
    function clearAddRouteStates() {
        if (lineStartMarker != null) {
            lineStartMarker.setAnimation(null); //跳动的动画
        }

        if (selectedMarker != null) {
            selectedMarker.setAnimation(null)
        }
        resetAddressEditBoxStatus(false)
    }

    function lineExist(points) {
        line = _.find(lines, function(l) { //排除包含了参数中所有点的线
            return _.every(points, function(p) {
                return _.contains(l.getPath(), p)
            })
        })
        if (line == null) {
            console.info("添加新的路径")
            return false
        } else {
            console.log("路径已经存在")
            return true
        }
    }

    function removeLine(points) {
        lines = _.reject(lines, function(l) { //排除包含了参数中所有点的线
            return _.every(points, function(p) {
                return _.contains(l.getPath(), p)
            })
        })
        console.log("剩余 %d 条路径", _.size(lines))
    }

    function removeMarker(bmapPoint) {
        markers = _.reject(markers, function(m) {
            return m.getPosition().equals(bmapPoint)
        })
        console.log("剩余 %d 个 Marker", _.size(markers))
    }

    function markerExist(bmapPoint) {
        var marker = _.find(markers, function(m) {
            // console.log(m.getPosition())
            // console.log(bmapPoint)
            return m.getPosition().equals(bmapPoint)
        })
        return marker != null
    }
    </script>
</body>

</html>
