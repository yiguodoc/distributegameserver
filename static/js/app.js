

require.config({
    paths: {
　　　　　　"jquery": "jquery.min",
　　　　　　"lodash": "lodash.min",
　　　　　　"Framework7": "../framework7/dist/js/framework7.min",
            "Chart": "Chart.min"
    },
    shim: {
        'Framework7': {
            deps: [],
            exports: 'Framework7'
        }
    }
});


require(['jquery', 'lodash', 'Framework7', 'Chart'], function ($, _, Framework7, Chart){

    // distributorID = $("body").attr("distributorID")
    // host = $("body").attr("host")
    // console.log("获取配送员编号：%s", distributorID)
    // Initialize your app
    myApp = new Framework7();
    // console.log(myApp.device.pixelRatio)
    // Export selectors engine
    $$ = Dom7;
    $$(document).on('pageAfterAnimation', function (e) {
      // Do something here when page loaded and initialized
        var page = e.detail.page
        console.log("%s loaded...", page.name)
        if(page.name == "processSelectOrder"){

        }
    })        
    mainView = myApp.addView('.view-main', {
        // Because we use fixed-through navbar we can enable dynamic navbar
        // dynamicNavbar: true
        domCache: true //enable inline pages
    });

    
    initMap()
    resetMap2Initial()

    mySwiper = swiper()
    pie()

    function swiper(){
        return myApp.swiper('.swiper-container', {
          pagination: '.swiper-pagination',
          paginationHide: false,
          paginationClickable: true,
          nextButton: '.swiper-button-next',
          prevButton: '.swiper-button-prev',
        }); 
    }

    // ==================================================
    // --------------------------------------------------
    //地图上的定位到当前控件
    // 定义一个控件类,即function
    function LocationControl(){
      // 默认停靠位置和偏移量
      this.defaultAnchor = BMAP_ANCHOR_TOP_LEFT;
      this.defaultOffset = new BMap.Size(10, 10);
    }

    // 通过JavaScript的prototype属性继承于BMap.Control
    LocationControl.prototype = new BMap.Control();

    // 自定义控件必须实现自己的initialize方法,并且将控件的DOM元素返回
    // 在本方法中创建个div元素作为控件的容器,并将其添加到地图容器中
    LocationControl.prototype.initialize = function(map){
      // 创建一个DOM元素
      var div = document.createElement("div");
      // 添加文字说明
      // div.appendChild(document.createTextNode("放大2级"));
      // 设置样式
      div.style.marginTop = "0px";
      div.style.width = "35px";
      div.style.height = "35px";
      div.style.cursor = "pointer";
      div.style.border = "0px solid gray";
      div.style.background = "url('/images/location.png') no-repeat 8px 8px";
      div.style.backgroundSize = "18px";
      div.style.backgroundColor = "white";
      div.style.opacity = "0.90";
      div.style.borderRadius = "3px";
      // 绑定事件,点击一次放大两级
      div.onclick = function(e){
        console.log("resetLocationButton clicked")
        if(myLocationMarker != null){
            var pos = myLocationMarker.getPosition()
            map.setCenter(pos)
        }
      }
      // 添加DOM元素到地图中
      map.getContainer().appendChild(div);
      // 将DOM元素返回
      return div;
    }
    // ==================================================

    // 创建控件添加到地图当中
    var myZoomCtrl = new LocationControl();
    map.addControl(myZoomCtrl);

    prepareConn()

    function prepareConn(){
        if (window["WebSocket"]) {
            conn = new WebSocket(wsUrl);
            conn.onclose = function(evt) {
                console.log("与服务器连接连接关闭，刷新重试")
            }
            conn.onopen = function(evt){
                console.log("与服务器连接成功")
            }
            conn.onmessage = function(evt) {
                var msg = evt.data
                // console.log(msg)
                msg = JSON.parse(msg)
                var handler = _.find(MessageHandlers, {MessageType: msg.MessageType})
                if(handler == null){
                    console.warn("消息处理函数未定义")
                }else{
                    handler.handler(msg)
                }
                // distributionProposals = JSON.parse(distributionProposals)
            }
        } else {
            alert("浏览器不支持 websocket")
        }        
    }

});

function t(){
    // Add view

    console.log("正在初始化基础数据")
    $$.get("/distributors?id={{.distributor.ID}}",function(data){
        console.log(data)
        if(_.isArray(data) == false){
            data = JSON.parse(data)
        }
        if(_.size(data) > 0){
            var distributor = data[0]
            console.log(distributor)
            console.log("游戏状态：%d", distributor.CheckPoint)
            if(distributor.CheckPoint <= 0){//还在初始化阶段
                viewRouteToPage(mainView, "process0")
            }else{//已经初始化过，中间可能掉线
                viewRouteToPage(mainView, "process1")
            }
        }
    })
}


//将地图设置为初始状态，目的是不突出任何信息
function resetMap2Initial(){

    // var point = new BMap.Point(116.644691, 39.934758);//北京物资学院
    // var point = new BMap.Point(116.331398,39.897445);//天安门
    setMapMarker(116.644691, 39.934758, false)

}
function initMap(){
    map = new BMap.Map("allmap");
    // map.centerAndZoom(new BMap.Point(116.404, 39.915), 16);
    map.enableScrollWheelZoom(true);
    var top_right_navigation = new BMap.NavigationControl({anchor: BMAP_ANCHOR_TOP_RIGHT, type: BMAP_NAVIGATION_CONTROL_ZOOM});//, offset: new BMap.Size(0, 40)
    map.addControl(top_right_navigation);    
    // map.addControl(new BMap.ZoomControl());  //添加地图缩放控件
    /*缩放控件type有四种类型:
    BMAP_NAVIGATION_CONTROL_SMALL：仅包含平移和缩放按钮；BMAP_NAVIGATION_CONTROL_PAN:仅包含平移按钮；BMAP_NAVIGATION_CONTROL_ZOOM：仅包含缩放按钮*/

}
function setMapMarker(lng,lat, bAddMarker){
    // map.removeOverlay(marker)
    if(lng > 0 || lat > 0){
        map.centerAndZoom(new BMap.Point(lng, lat), 16);
        // map.addControl(new BMap.ZoomControl());  //添加地图缩放控件
        if(bAddMarker == true){
            var marker = new BMap.Marker(new BMap.Point(lng, lat));  //创建标注
            map.addOverlay(marker);                 // 将标注添加到地图中            
        }
    }
}
// initMap()
// resetMap2Initial()
function viewRouteToPage(view, page){

    // mainView.router.load({pageName:page})
    view.router.load({pageName:page})
}