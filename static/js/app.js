
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
// initMap()
// resetMap2Initial()
function viewRouteToPage(view, page){

    // mainView.router.load({pageName:page})
    view.router.load({pageName:page})
}