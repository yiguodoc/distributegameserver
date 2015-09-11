

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
    // pie()

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
                myApp.alert("您已掉线，5秒后重新连接", "提示")
                setTimeout(function () {
                        myApp.closeModal();
                        // prepareConn()
                        window.location.href = window.location.href
                    }, 5000);
            }
            conn.onopen = function(evt){
                console.log("与服务器连接成功")
            }
            conn.onmessage = function(evt) {
                var msg = evt.data
                msg = JSON.parse(msg)
                if(msg.ErrorMsg != null && msg.ErrorMsg.length > 0){//系统对操作的错误提示信息
                    console.warn(msg.ErrorMsg)
                    addMsgToView(msg.SysTime, msg.ErrorMsg)
                }
                var handler = _.find(MessageHandlers, {MessageType: msg.MessageType})
                if(handler == null){
                    console.warn("消息处理函数未定义")
                }else{
                    if(handler.print != false){
                        console.log("系统消息：",msg)
                    }
                    handler.handler(msg)
                }
                // distributionProposals = JSON.parse(distributionProposals)
            }
        } else {
            alert("浏览器不支持 websocket")
        }        
    }

});

function getRad(d){   
    var PI = Math.PI;    
   return d*PI/180.0;    
}     
/** 
     * approx distance between two points on earth ellipsoid ,return the distance  in float meter.
     * @param {Object} lat1   
     * @param {Object} lng1   
     * @param {Object} lat2   
     * @param {Object} lng2   
     */     
function CoolWPDistance(lat1,lng1,lat2,lng2){     
    var f = getRad((lat1 + lat2)/2);     
    var g = getRad((lat1 - lat2)/2);     
    var l = getRad((lng1 - lng2)/2);     
    var sg = Math.sin(g);     
    var sl = Math.sin(l);     
    var sf = Math.sin(f);     
    var s,c,w,r,d,h1,h2;     
    var a = 6378137.0;//The Radius of eath in meter.   
    var fl = 1/298.257;     
    sg = sg*sg;     
    sl = sl*sl;     
    sf = sf*sf;     
    s = sg*(1-sl) + (1-sf)*sl;     
    c = (1-sg)*(1-sl) + sf*sl;     
    w = Math.atan(Math.sqrt(s/c));     
    r = Math.sqrt(s*c)/w;     
    d = 2*w*a;     
    h1 = (3*r -1)/2/c;     
    h2 = (3*r +1)/2/s;     
    s = d*(1 + fl*(h1*sf*(1-sg) - h2*(1-sf)*sg));   
    // s = s/1000;   
    // s = s.toFixed(2);//指定小数点后的位数。   
    return s;     
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
    console.info("page => %s",page)
    // mainView.router.load({pageName:page})
    view.router.load({pageName:page})
}