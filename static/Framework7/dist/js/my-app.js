// Initialize your app
var myApp = new Framework7();

// Export selectors engine
var $$ = Dom7;

// Add view
var mainView = myApp.addView('.view-main', {
    // Because we use fixed-through navbar we can enable dynamic navbar
    dynamicNavbar: true
});
var marker = null
var map = new BMap.Map("allmap");



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
initMap()
resetMap2Initial()

// mainView.Router.load("index")
// flashView.Router.load("index")
// Callbacks to run specific code for specific pages, for example for About page:
myApp.onPageInit('about', function (page) {
    // run createContentPage func after link was clicked
    $$('.create-page').on('click', function () {
        createContentPage();
    });
});

// Generate dynamic page
var dynamicPageIndex = 0;
function createContentPage() {
	mainView.router.loadContent(
        '<!-- Top Navbar-->' +
        '<div class="navbar">' +
        '  <div class="navbar-inner">' +
        '    <div class="left"><a href="#" class="back link"><i class="icon icon-back"></i><span>Back</span></a></div>' +
        '    <div class="center sliding">Dynamic Page ' + (++dynamicPageIndex) + '</div>' +
        '  </div>' +
        '</div>' +
        '<div class="pages">' +
        '  <!-- Page, data-page contains page name-->' +
        '  <div data-page="dynamic-pages" class="page">' +
        '    <!-- Scrollable page content-->' +
        '    <div class="page-content">' +
        '      <div class="content-block">' +
        '        <div class="content-block-inner">' +
        '          <p>Here is a dynamic page created on ' + new Date() + ' !</p>' +
        '          <p>Go <a href="#" class="back">back</a> or go to <a href="services.html">Services</a>.</p>' +
        '        </div>' +
        '      </div>' +
        '    </div>' +
        '  </div>' +
        '</div>'
    );
	return;
}