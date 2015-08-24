<!DOCTYPE html>
<html>
  <head>
    <!-- Required meta tags-->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no, minimal-ui">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black">
    <!-- Your app title -->
    <title>配送大师</title>
    <!-- Path to Framework7 iOS CSS theme styles-->
    <link rel="stylesheet" href="framework7/dist/css/framework7.ios.min.css">
    <!-- Path to Framework7 iOS related color styles -->
    <link rel="stylesheet" href="framework7/dist/css/framework7.ios.colors.min.css">
    <!-- Path to your custom app styles-->
    <link rel="stylesheet" href="stylesheets/app.css">


  </head>
  <body class="theme-lightblue" distributorID = "{{.distributor.ID}}">
    <!-- Status bar overlay for full screen mode (PhoneGap) -->
    <div class="statusbar-overlay"></div>
    <!-- Panels overlay-->
    <div class="panel-overlay"></div>
    <!-- Left panel, let it be with reveal effect -->
    <div class="panel panel-left panel-reveal">
        左侧栏
    </div>
    <!-- Views -->
    <!-- <div class="views"> -->
    <div class="views tabs toolbar-through">

		<!-- Your main view, should have "view-main" class -->
		<!-- 信息页面 -->
		<div class="view view-main tab active" id="view-main">
			<div class="navbar">
			    <div class="navbar-inner">
			      <div class="center sliding">配送大师</div>
			      <div class="right">
			          <a href="#" class="link icon-only open-panel">
			              <i class="icon icon-bars"></i>
			          </a>
			      </div>
			    </div>
			</div>
		      <!-- Pages container, because we use fixed-through navbar and toolbar, it has additional appropriate classes-->
		    <div class="pages navbar-through">
		        <!-- Page, "data-page" contains page name -->
		        
		        <div data-page="process1" class="page" id="1">
		              <!-- Scrollable page content -->
		              <div class="page-content "> 
				        <div class="content-block" style="margin-top: 20px;  margin-bottom: 15px;">
			                <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
			                <div id="canvas-holder" style="text-align: center;">
	                			<canvas id="chart-area" width="130" height="130"/>
	                		</div>
			                <p style="text-align: center;margin-top:0px;">订单区域分布比例</p>

			            </div>
			            <div class="swiper-custom">
			               <div class="swiper-container">
			                 <div class="swiper-pagination"></div>
			                 <div class="swiper-wrapper">
			                   <div class="swiper-slide">
				                   	<span class="slide-title">订单编号01</span>
				                   	<span class="slide-content">地址01</span>
			                   </div>

			                   <div class="swiper-slide">
				                   	<span class="slide-title">订单编号02</span>
				                   	<span class="slide-content">地址02</span>
			                   </div>

			                   <div class="swiper-slide">
				                   	<span class="slide-title">订单编号03</span>
				                   	<span class="slide-content">地址03</span>
			                   </div>
			                 </div>
			               </div>
			               <div class="swiper-button-prev"></div>
			               <div class="swiper-button-next"></div>
			             </div>

		                <div class=" login-btn-content">
		                	<div class="row">
		                		<div class="col-10"></div>
		                		<div class="col-80">
				                     <a href="#" class="button button-big button-fill" id="" onclick="selectOrder()">选择订单</a>
		                		</div>
		                		<div class="col-10"></div>
		                    </div>
		                </div>

		              </div>
		        </div>
                <div data-page="waiting" class="page" id="2">
                      <!-- Scrollable page content -->
                      <div class="page-content "> 
        		        <div class="content-block" style="margin-top: 100px;">
        		                <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
        		                <p style="text-align: center;">等待中...</p>
        		                <div class=" login-btn-content">
        		                      <a href="#process1" class="button button-big button-fill" id="" onclick="viewRouteToPage(mainView, 'process1')">进入游戏</a>
        		                </div>

        	            </div>

                      </div>
                 </div>

		        <div data-page="index" class="page" id="3">
		              <!-- Scrollable page content -->
		              <div class="page-content "> 
				        <div class="content-block" style="margin-top: 100px;">
				                <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
				                <p style="text-align: center;">我准备好了</p>
				                <div class=" login-btn-content">
				                      <a href="#" class="button button-big button-fill" id="" onclick="onPreparedToStartGame()">进入游戏</a>
				                </div>
			            </div>

		              </div>
		         </div>


		    </div>
		</div>
		
		<!-- 地图页面 -->
		<div class="view tab" id="view-map">
			<div class="navbar">
			    <div class="navbar-inner">
			      <div class="center sliding">配送大师</div>
			      <div class="right">
			          <a href="#" class="link icon-only open-panel">
			              <i class="icon icon-bars"></i>
			          </a>
			      </div>
			    </div>
			</div>
		      <!-- Pages container, because we use fixed-through navbar and toolbar, it has additional appropriate classes-->
		    <div class="pages navbar-through toolbar-through">
		        <!-- Page, "data-page" contains page name -->
		        <div data-page="map" class="page" id="pageMap">
		              <!-- Scrollable page content -->
		              <div class="page-content "> 



				        <div id="allmap" style="height:99%;margin-top:1px;">

				        </div>
			        	<div class="row" style="top: 55px; position: absolute; height: 35px;left:43px;right:45px; opacity: 0.8; background-color: white;  border-left: 1px solid gray; border-right: 1px solid gray;">
								<div style="width:100%; text-align: center; color: gray; margin-top: 10px;">北京物资学院</div>
			        		<div class="col-10"> </div>

			        		<div class="col-80"> 
			        		</div>
			        		<div class="col-10"> </div>
			        	</div>
			        	<div class="row" style="margin-top:-80px;">
			        		<div class="col-10"> </div>
			        		<div class="col-20">
								<a href="#" class="button button-big button-fill" id="" onclick="">&lt;</a>
			        		</div>
			        		<div class="col-40"> 
								<a href="#" class="button button-big button-fill color-lightblue" id="" onclick="">去往该点</a>
			        		</div>
			        		<div class="col-20">
								<a href="#" class="button button-big button-fill" id="" onclick="">&gt;</a>
			        		 </div>
			        		<div class="col-10"> </div>
			        	</div>
							
		              </div>
		        
		          </div>
		    </div>
		</div>		

		<!-- 订单页面 -->
		<div class="view tab" id="view-orders">
			<div class="navbar">
			    <div class="navbar-inner">
			      <div class="center sliding">配送大师</div>
			      <div class="right">
			          <a href="#" class="link icon-only open-panel">
			              <i class="icon icon-bars"></i>
			          </a>
			      </div>
			    </div>
			</div>
		      <!-- Pages container, because we use fixed-through navbar and toolbar, it has additional appropriate classes-->
		    <div class="pages navbar-through toolbar-through">
		        <!-- Page, "data-page" contains page name -->
		        <div data-page="map" class="page" id="pageMap">
		              <!-- Scrollable page content -->
		              <div class="page-content "> 
				        订单页面

		              </div>
		         

		          </div>
		    </div>
		</div>

		<!-- 道具页面 -->
		<div class="view tab" id="view-cards">
			<div class="navbar">
			    <div class="navbar-inner">
			      <div class="center sliding">配送大师</div>
			      <div class="right">
			          <a href="#" class="link icon-only open-panel">
			              <i class="icon icon-bars"></i>
			          </a>
			      </div>
			    </div>
			</div>
		      <!-- Pages container, because we use fixed-through navbar and toolbar, it has additional appropriate classes-->
		    <div class="pages navbar-through toolbar-through">
		        <!-- Page, "data-page" contains page name -->
		        <div data-page="map" class="page" id="pageMap">
		              <!-- Scrollable page content -->
		              <div class="page-content "> 
				        道具页面

		              </div>
		         

		          </div>
		    </div>
		</div>


		<div class="toolbar tabbar tabbar-labels">
		  <div class="toolbar-inner">
		      <a href="#view-main" class="tab-link active">
		          <i class="icon tabbar-demo-icon-1"></i>
		          <span class="tabbar-label">状态</span>
		      </a>
		      <a href="#view-map" class="tab-link ">
		          <i class="icon tabbar-demo-icon-2">
		              <span class="badge bg-red">5</span>
		          </i>
		          <span class="tabbar-label">地图</span>
		      </a>
		      <a href="#view-orders" class="tab-link">
		          <i class="icon tabbar-demo-icon-3"></i>
		          <span class="tabbar-label">订单</span>
		      </a>
		      <a href="#view-cards" class="tab-link">
		          <i class="icon tabbar-demo-icon-4"></i>
		          <span class="tabbar-label">消息</span>
		      </a>
		  </div>
		</div> 


    </div>
    <script type="text/javascript" src="http://api.map.baidu.com/api?v=2.0&ak=kU4NWwyP5SwguC2W2WAfO1bO"></script>

    <script type="text/javascript">
	    var marker = null
	    var conn;
	    var distributorID = "{{.distributor.ID}}"
	    var orders = []
	    var mainView
	    var map
	    var $$
	    var myApp
	    var mySwiper 
	    var wsUrl = "ws://{{.HOST}}/wsOrderDistribution?id={{.distributor.ID}}" 
	    var MessageHandlers = [
	    	{MessageType: {{.pro_2c_all_prepared_4_select_order}}, handler: function(msg){
	    		console.log("route to %s", 'process1')
	        	viewRouteToPage(mainView, 'process1')
	    	}},
	    	{MessageType: {{.pro_2c_order_distribution_proposal}}, handler: pro_2c_order_distribution_proposal_handler},
	    	{MessageType: {{.pro_timer_count_down}}, handler: pro_timer_count_down_handler},
	    	{MessageType: {{.pro_2c_message_broadcast}}, handler: pro_2c_message_broadcast_handler},
	    	{MessageType: {{.pro_2c_order_select_result}}, handler: pro_2c_order_select_result_handler},
	    	{MessageType: {{.pro_2c_order_full}}, handler: pro_2c_order_full_handler},
	    	{MessageType: {{.pro_2c_distributor_info}}, handler: pro_2c_distributor_info_handler}
	    ]


	    function pro_2c_order_distribution_proposal_handler(msg){
	    	mySwiper.removeAllSlides();
	    	orders = msg.Data
	    	_.each(orders, function(order){
	    	    var orderTip = "编号："+order.ID
	    	    if (order.GeoSrc != null){
	    	        orderTip += "  位置:"+order.GeoSrc.Address
	    	    } 
	    	    mySwiper.appendSlide(String.format('<div class="swiper-slide"> <span class="slide-title">{0}</span> <span class="slide-content">{1}</span> </div>',order.ID, orderTip))
	    	})
	    }
	    function pro_timer_count_down_handler(msg){
            console.log("-> "+ msg.Data)
            if(msg.Data <= 1){
	        	viewRouteToPage(mainView, 'process1')
	        	console.log("route to process1")
            }	    	
	    }
	    function pro_2c_message_broadcast_handler(msg){
            console.log(msg.Data)
	    }
	    function pro_2c_order_select_result_handler(msg){
	    	if(msg.Data.DistributorID == distributorID){
	    	    console.log("抢到了订单 "+msg.Data.OrderID)
	    	}else{
	    	    console.log("没有抢到订单 "+msg.Data.OrderID)
	    	}
	    	console.log(msg.Data)
	    }
	    function pro_2c_order_full_handler(msg){
	    	if(msg.Data == distributorID){
	    	    hideOrderSelectButton()
	    	    $("#btnStartDistribute").show()
	    	}
	    }
	    function pro_2c_distributor_info_handler(msg){
	    	var data = msg.Data
	    	console.log(data)
	    	distributor = data
	    	if(distributor.CheckPoint <= 0){//还在初始化阶段
	    	    // $("#btnPrepared").show()
	    	}else{//已经初始化过，中间可能掉线
	    	    // showOrderSelectButton()
	    	}
	    }

	    function onPreparedToStartGame(){
        	viewRouteToPage(mainView, 'waiting')
	        	// viewRouteToPage(mainView, 'process1')
        	send({MessageType: {{.pro_prepared_for_select_order}}, Data:{DistributorID: distributorID}})
	    }
	    function send(msg){
	        if (!conn) {
	            return false;
	        }
	        conn.send(JSON.stringify(msg))
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
    </script>
	<script src="javascripts/require.js" data-main="javascripts/app"></script>

  </body>
</html>    