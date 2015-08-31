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
		        <!-- 订单选择页面 -->
		        <div data-page="processSelectOrder" class="page" id="1">
		              <div class="page-content "> 
				        <div class="content-block" style="margin-top: 20px;  margin-bottom: 15px;">
			                <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
			                <div id="canvas-holder" style="text-align: center;">
	                			<canvas id="chart-area" width="130" height="130"/>
	                		</div>
			                <p style="text-align: center;margin-top:0px;font-size: 12px;">订单区域分布比例</p>

			            </div>
			            <div class="swiper-custom">
			               <div class="swiper-container">
			                 <div class="swiper-pagination"></div>
			                 <div class="swiper-wrapper">
								<!--<div class="swiper-slide">
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
			                   </div> -->
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
		        <!-- 等待其他参与者进入的页面 -->
                <div data-page="waiting" class="page" id="2">
                      <!-- Scrollable page content -->
                      <div class="page-content "> 
        		        <div class="content-block" style="margin-top: 100px;">
        		                <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
        		                <p id="waitingInfo" style="text-align: center;">等待其他人进入...</p>
        		                <!-- <div class=" login-btn-content">
        		                      <a href="#processSelectOrder" class="button button-big button-fill" id="" onclick="viewRouteToPage(mainView, 'processSelectOrder')">进入游戏</a>
        		                </div> -->

        	            </div>

                      </div>
                </div>
                <!-- 订单选择完毕，转入配送状态页面之前的页面 -->
                <div data-page="processGo2Distribution" class="page" id="4">
					<div class="page-content "> 
					    <div class="content-block" style="margin-top: 100px;">
					            <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
					            <p style="text-align: center;">订单选择已经完成</p>
					            <div class=" login-btn-content">
					                  <a href="#" class="button button-big button-fill" id="" onclick="onPreparedToStartGame()">开始配送</a>
					            </div>
					    </div>

					</div>
                </div>
                <!-- 配送状态页面 -->
                <div data-page="processDistribution" class="page" id="5"><!-- Scrollable page content -->
                       
					<div class="page-content "> 
					    <div class="content-block" style="margin-top: 0px;">
					            <!-- <p style="text-align: center;">00:10</p> -->
					            <div class="" style="text-align: center;margin-top:5px;  color: rgba(150,150,150,0.9);">00:10</div>
					            <div class="content-block-inner" style="text-align: center;margin-top:50px;">
					            	<div style="font-size: 13px; color: rgba(150,150,150,0.8);">订单签收进度</div>
					            	<div style="font-size: 19px;">7/10</div>

					            </div>
					            <div style = "  text-align: center; margin-bottom: 40px; font-size: 28px; border-bottom: 1px solid rgba(200,200,200,0.5); padding-bottom: 60px; padding-top: 60px;">
					            	<span>配送中</span>
					            </div>
					            <div class=" login-btn-content">
					                  <a href="#" class="button button-big button-fill disabled" id="btnSignOrder" onclick="onSignOrder()">订单签收</a>
					            </div>
					    </div>

					</div>
                </div>
                <!-- 配送完成统计页面 -->
                <div data-page="processStatistic" class="page" id="6"><!-- Scrollable page content -->
                       
					<div class="page-content "> 
					    <div class="content-block" style="margin-top: 0px;">
					            <!-- <p style="text-align: center;">00:10</p> -->
					            <div class="" style="text-align: center; margin-top: 55px; color: rgba(100,100,100,1.9); font-size: 20px;">订单配送已完成</div>
					            <div style = "text-align: center; font-size: 15px; border-bottom: 1px solid rgba(200,200,200,0.4); padding-top: 40px; color: rgba(100,100,100,0.5);">
					            	<span>成绩统计</span>
					            </div>
					            <div class="content-block-inner" style="margin-top:0px;">
					            	<div class="row no-gutter" style="margin-top:5px;">
					            		<div class="col-50" style="text-align: right;padding-right: 3px;"> 排名 </div>
					            		<div class="col-50" style="text-align: left;padding-left: 3px;"> 2</div>
					            	</div>
					            	<div class="row no-gutter" style="margin-top:5px;">
					            		<div class="col-50" style="text-align: right;padding-right: 3px;"> 总耗时 </div>
					            		<div class="col-50" style="text-align: left;padding-left: 3px;"> 17分20秒</div>
					            	</div>


<!-- 					            	<div style="font-size: 13px; color: rgba(150,150,150,0.8);">
					            		<span>总耗时：</span> <span>17分20秒</span>
					            	</div>
					            	<div style="font-size: 19px;">
					            		<span>排名：</span> <span>2</span>
					            	</div> -->

					            </div>
					            
					    </div>

					</div>
                </div>

                <!-- 登录进入游戏之后的页面 -->
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
			        	<div class="row" style="margin-top:-60px;">
			        		<div class="col-10"> </div>
			        		<div class="col-20">
								<a href="#" class="button button-fill" id="" onclick="onFocusOnPreNode()">&lt;</a>
			        		</div>
			        		<div class="col-40"> 
								<a href="#" class="button  button-fill color-lightblue" id="" onclick="onChooseDestNode()">去往该点</a>
			        		</div>
			        		<div class="col-20">
								<a href="#" class="button  button-fill" id="" onclick="onFocusOnNextNode()">&gt;</a>
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
				        <div  class="content-block-title" style="margin-top: 10px;">未签收<span id="orderListUnsignedLabel">（0）</span></div>
				        <div class="list-block media-list">
				          <ul id="orderListUnsigned">
				            <li>
				              <a href="#" class="item-link item-content">
				                <div class="item-inner">
				                  <div class="item-title-row">
				                    <div class="item-title">Facebook</div>
				                    <div class="item-after"><span class="badge">+1</span></div>
				                  </div>
				                  <!-- <div class="item-subtitle">New messages from John Doe</div> -->
				                  <div class="item-text">Lorem ipsum dolor sit amet...</div>
				                </div>
				              </a>
				            </li>
				            
				          </ul>
				        </div>

				        <div class="content-block-title" style="margin-top: 10px;">已签收<span  id="orderListSignedLabel">（0）</span></div>
				        <div   class="list-block media-list">
				          <ul id="orderListSigned">
				            <li>
				              <a href="#" class="item-link item-content">
				                <div class="item-inner">
				                  <div class="item-title-row">
				                    <div class="item-title">Facebook</div>
				                    <div class="item-after">17:14</div>
				                  </div>
				                  <!-- <div class="item-subtitle">New messages from John Doe</div> -->
				                  <div class="item-text">Lorem ipsum dolor sit amet...</div>
				                </div>
				              </a>
				            </li>
				            
				            <li>
				              <a href="#" class="item-link item-content">
				                <div class="item-inner">
				                  <div class="item-title-row">
				                    <div class="item-title">Facebook</div>
				                    <div class="item-after">17:14</div>
				                  </div>
				                  <!-- <div class="item-subtitle">New messages from John Doe</div> -->
				                  <div class="item-text">Lorem ipsum dolor sit amet...</div>
				                </div>
				              </a>
				            </li>
				            
				          </ul>
				        </div>

		              </div>
		         

		          </div>
		    </div>
		</div>

		<!-- 消息页面 -->
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
				         <div class="list-block media-list" style="margin-top: 0px; padding-left: 5px;">
				           <ul id="msgList">
				           <!--   <li>
				                 <div class="item-inner">
				                   <div class="item-title-row">
				                     <div class="item-title">12:34</div>
				                     <div class="item-after"></div>
				                   </div>
				                   <div class="item-text">New messages from John DoeNew messages from John DoeNew messages from John Doe</div>
				                 </div>
				             </li>
				             
				             <li>
				                 <div class="item-inner">
				                   <div class="item-title-row">
				                     <div class="item-title">12:34</div>
				                     <div class="item-after"></div>
				                   </div>
				                   <div class="item-text">Lorem ipsum dolor sit amet...</div>
				                 </div>
				             </li> -->
				             
				           </ul>
				         </div>

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
	    var mainView, map, mapData, $$, myApp, mySwiper, conn;
	    var distributorID = "{{.distributor.ID}}"
	    // var orders = []
	    var wsUrl = "ws://{{.HOST}}/wsOrderDistribution?id={{.distributor.ID}}" 
	    var myLocationMarker = null
	    var markerAim = null//准星环绕的marker
	    var markerDest = null//目的地的准星
	    var nextNodeSelector = null
		var iconKinds = [
    		{type:0, imageName: "warehouse.png", width: 64, height: 64, opt: {anchor: new BMap.Size(32, 48)}},
    		{type:1, imageName: "nodeSmall.png", width: 12, height: 12, opt: {anchor: new BMap.Size(6, 6)}},
    		{type:2, imageName: "nodeSmall.png", width: 12, height: 12, opt: {anchor: new BMap.Size(6, 6)}},
    		{hasOrder: true, orderSigned: false, imageName: "bagageClosed.png", width: 29, height: 29, opt: {anchor: new BMap.Size(15, 15)}},
    		{hasOrder: true, orderSigned: true, imageName: "bagageOpen.png", width: 29, height: 29, opt: {anchor: new BMap.Size(15, 15)}},
    		{canBeSelected: true, imageName: "aimBlack.png", width: 100, height: 100, opt: {anchor: new BMap.Size(16, 16), imageSize: new BMap.Size(32,32)}},
    		{selected: true, imageName: "aimRed.png", width: 100, height: 100, opt: {anchor: new BMap.Size(16, 16), imageSize: new BMap.Size(32,32)}},
    		{}
		]
	    var MessageHandlers = [
	    	{MessageType: {{.pro_2c_all_prepared_4_select_order}}, handler: function(msg){
	    		console.log("route to %s", 'processSelectOrder')
	        	viewRouteToPage(mainView, 'processSelectOrder')
	    	}},
	    	{MessageType: {{.pro_2c_message_broadcast_before_game_start}}, handler: pro_2c_message_broadcast_before_game_start_handler},
	    	{MessageType: {{.pro_2c_order_distribution_proposal}}, handler: pro_2c_order_distribution_proposal_handler},
	    	{MessageType: {{.pro_timer_count_down}}, handler: pro_timer_count_down_handler},
	    	{MessageType: {{.pro_2c_message_broadcast}}, handler: pro_2c_message_broadcast_handler},
	    	{MessageType: {{.pro_2c_order_select_result}}, handler: pro_2c_order_select_result_handler},
	    	{MessageType: {{.pro_2c_order_full}}, handler: pro_2c_order_full_handler},
	    	{MessageType: {{.pro_2c_reach_route_node}}, handler: pro_2c_reach_route_node_handler},
	    	{MessageType: {{.pro_2c_move_to_new_position}}, handler: pro_2c_move_to_new_position_handler},
	    	{MessageType: {{.pro_2c_move_from_node}}, handler: pro_2c_move_from_node_handler},
	    	{MessageType: {{.pro_2c_map_data}}, handler: pro_2c_map_data_handler},
	    	{MessageType: {{.pro_2c_distributor_info}}, handler: pro_2c_distributor_info_handler},
	    	{MessageType: {{.pro_2c_reset_destination}}, handler: pro_2c_reset_destination_handler},
	    	{MessageType: {{.pro_2c_sign_order}}, handler: pro_2c_sign_order_handler},
	    	{MessageType: {{.pro_2c_all_order_signed}}, handler: pro_2c_all_order_signed_handler},
	    	{}
	    ]
	    function pro_2c_all_order_signed_handler(msg){
			distributor = msg.Data
        	viewRouteToPage(mainView, 'processStatistic')

	    }
	    function pro_2c_sign_order_handler(msg){
	    	if(msg.Data == null){
		    	console.info("订单签收失败")
	    		return
	    	}else{
	    		distributor = msg.Data
		    	console.info("订单签收成功")
		    	//按照已经签收的订单的位置将图标重置为订单接收后的状态图标
		    	_.each(distributor.AcceptedOrders, function(order){
		    		if(order.Signed == true){
			    		var pos = order.GeoSrc
			    		var pos = _.findWhere(mapData.Points, {Lat: pos.Lat, Lng: pos.Lng})
			    		if(pos != null){
				    		addPointMarkerToMap(pos, {hasOrder: true, orderSigned: true}, true)
			    		}else{
			    			console.warn("订单位置未找到")
			    		}
		    		}
		    	})
		    	remindOrderSigning()
				refreshOderView()
	    	}
	    }
	    function pro_2c_reset_destination_handler(msg){
	    	console.info("重置目的地成功")
    		distributor = msg.Data
	    	setDestinationMarker()
			refreshNodeToSelect()
	    }
	    function pro_2c_reach_route_node_handler(msg){
    		distributor = msg.Data
			refreshMyLocation()
			console.info("配送员到达节点")
			refreshNodeToSelect()
	    	setDestinationMarker()
	    	remindOrderSigning()
	    }

	    function pro_2c_move_to_new_position_handler(msg){
    		distributor = msg.Data
			refreshMyLocation()
			console.info("配送员位置发生变化")
	    }
	    function pro_2c_move_from_node_handler(msg){
    		distributor = msg.Data
			refreshMyLocation()
			console.info("配送员离开节点")
			refreshNodeToSelect()
	    	resetOrderSignButtonState(false)	    	
	    	// setDestinationMarker()
	    }
	    function pro_2c_message_broadcast_before_game_start_handler(msg){
	    	$$("#waitingInfo").text(msg.Data)
	    }
	    function pro_2c_order_distribution_proposal_handler(msg){
	    	mySwiper.removeAllSlides();
	    	var orders = msg.Data
	    	_.each(orders, function(order){
	    	    mySwiper.appendSlide(String.format('<div class="swiper-slide"> <span class="slide-title" style="background-color: rgba({0}, 0.6);">{1}</span> <span class="slide-content">{2}北京市物资学院</span> </div>',order.Region.Color, order.ID, order.GeoSrc.Address))
	    	})
	    }

	    function pro_timer_count_down_handler(msg){
            console.log("-> "+ msg.Data)	    	
	    }
	    function pro_2c_message_broadcast_handler(msg){
            console.log(msg.Data)
	    }
	    function pro_2c_order_select_result_handler(msg){
	    	console.log(msg.Data)
	    	if(msg.Data != null){
	    	    console.log("抢到了订单 ")
	    		distributor = msg.Data
	    		resetPie()
				flagOrderNodeMarker()
				refreshOderView()
				var lastOrder = distributor.AcceptedOrders[_.size(distributor.AcceptedOrders)-1]
				addMsgToView({timeStamp: transformTimeElapseToStandardFormat(lastOrder.SelectedTime), content: "选择配送订单 "+ lastOrder.ID})
	    	}else{
	    	    console.log("没有抢到订单 ")	    		
	    	}
	    }
	    function pro_2c_order_full_handler(msg){
	    	distributor = msg.Data
        	viewRouteToPage(mainView, 'processDistribution')
	    }
	    function pro_2c_distributor_info_handler(msg){
	    	var data = msg.Data
	    	console.log(data)
	    	distributor = data
	    	switch(distributor.CheckPoint){
	    		case {{.checkpoint_flag_origin}}:
		        	viewRouteToPage(mainView, 'index')
	    		break
	    		case {{.checkpoint_flag_order_select}} :
		        	viewRouteToPage(mainView, 'processSelectOrder')
		        	resetPie()
	    		break
	    		case {{.checkpoint_flag_order_distribute}} :
		        	viewRouteToPage(mainView, 'processDistribution')
		        	// viewRouteToPage(mainView, 'processGo2Distribution')
	    		break
	    		case {{.checkpoint_flag_order_distribute_over}}:
		        	viewRouteToPage(mainView, 'processStatistic')
	    		break
	    	}
			refreshMyLocation()
			flagOrderNodeMarker()
			refreshNodeToSelect()
			remindOrderSigning()
			refreshOderView()
			addMsgToView({timeStamp: transformTimeElapseToStandardFormat(distributor.TimeElapse), content: "上线"})
	    }
	    function pro_2c_map_data_handler(msg){
	    	if(mapData != null){
	    		_.each(mapData.Points, function(point){
	    			map.removeOverlay(point.mark)
	    		})
	    		_.each(mapData.Lines, function(line){
	    			map.removeOverlay(line.lineOverlay)
	    		})
	    	}else{
		    	// map.clearOverlays()
	    	}
	    	mapData = msg.Data
	    	drawRouteNodeOnMap(mapData)
	    }
	    //---------------------------------------------------------
	    //添加消息到消息页面
	    function addMsgToView(msg){
	    	var dom = $$("#msgList")
	    	dom.prepend(String.format('<li> <div class="item-inner"> <div class="item-title-row"> <div class="item-title">{0}</div> <div class="item-after"></div> </div>  <div class="item-text">{1}</div> </div> </li>', msg.timeStamp, msg.content))
	    }

	    //重置订单页面的订单列表，当订单数量或者状态发生变化时调用
	    function refreshOderView(){
	    	var refresher = function(domID, orders){
	    		var dom = $$("#"+domID)
		    	dom.children().remove()

		    	$$("#"+domID+"Label").text(String.format("（{0}）", _.size(orders)))
	    		_.each(orders, function(order){
	    			dom.append(String.format('<li> <a href="#" class="item-link item-content"> <div class="item-inner"> <div class="item-title-row"> <div class="item-title">{0}</div> <div class="item-after"><span class="badge  bg-lightblue">+1</span></div> </div> <div class="item-text">{1}</div> </div> </a> </li>', order.ID, order.Address || ""))
	    		})
	    	}
	    	refresher("orderListUnsigned", _.filter(distributor.AcceptedOrders, function(order){return order.Signed == false}))
	    	refresher("orderListSigned", _.filter(distributor.AcceptedOrders, function(order){return order.Signed == true}))

	    	// var orderListUnsigned = $$("#orderListUnsigned")
	    	// orderListUnsigned.children().remove()
	    	// var ordersUnsigned = _.filter(distributor.AcceptedOrders, function(order){return order.Signed == false})
	    	// _.each(ordersUnsigned, function(order){
	    	// 	orderListUnsigned.append(String.format('<li> <a href="#" class="item-link item-content"> <div class="item-inner"> <div class="item-title-row"> <div class="item-title">{0}</div> <div class="item-after">17:14</div> </div> <div class="item-text">{1}</div> </div> </a> </li>', order.ID, order.Address))
	    	// })

	    	// var orderListSigned = $$("orderListSigned")
	    	// orderListSigned.children().remove()
	    	// var ordersSigned = _.filter(distributor.AcceptedOrders, function(order){return order.Signed == true})
	    	// _.each(ordersSigned, function(order){
	    	// 	orderListSigned.append(String.format('<li> <a href="#" class="item-link item-content"> <div class="item-inner"> <div class="item-title-row"> <div class="item-title">{0}</div> <div class="item-after">17:14</div> </div> <div class="item-text">{1}</div> </div> </a> </li>', order.ID, order.Address))
	    	// })
	    }
	    function remindOrderSigning(){
	    	//如果有订单，提醒签收
	    	var pos = distributor.CurrentPos
	    	var orderTemp = _.find(distributor.AcceptedOrders, function(order){
	    		return order.GeoSrc.Lat == pos.Lat && order.GeoSrc.Lng == pos.Lng && order.Signed == false
	    	})
	    	if(orderTemp != null){
	    		myApp.alert('有订单，注意签收', '提醒');
	    		setTimeout(function () {
	    		    myApp.closeModal()
	    		}, 2000);
		    	resetOrderSignButtonState(true)	    	
	    	}else{
		    	resetOrderSignButtonState(false)	    	

	    	}
	    }
	    function pointsAimAtLooper(points){
	    	this.points= points
	    	this.currentIndex = 0
	    	this.removeLast = function(){
		    	if(markerAim != null){
		    		map.removeOverlay(markerAim)
		    		markerAim = null
		    	}
		    	return this
	    	}
	    	this.removeLast()
	    	this.setMarker = function(){
	    		this.removeLast()
	    		if(this.currentIndex >= 0 && (_.size(this.points) > this.currentIndex)){
		    		markerAim = addPointMarkerToMap(this.points[this.currentIndex], {canBeSelected: true}, false)
	    		}	    		
	    	}
	    	this.next = function(){
	    		this.currentIndex ++
	    		if(this.currentIndex >= _.size(this.points)){
	    			this.currentIndex = 0
	    		}
		    	this.setMarker()
	    	}
	    	this.pre = function(){
	    		this.currentIndex --
	    		if(this.currentIndex < 0){
	    			this.currentIndex = _.size(this.points) - 1
	    		}
	    		this.setMarker()
	    	}
	    	this.currentPostion = function(){
	    		if(_.size(this.points) > 0){
		    		return this.points[this.currentIndex]
	    		}else{
	    			return null
	    		}
	    	}
	    }
	    function setDestinationMarker(){
	    	if(markerDest != null){
	    		map.removeOverlay(markerDest)
	    		markerDest = null
	    	}
	    	if(distributor.DestPos != null){
	    		markerDest = addPointMarkerToMap(distributor.DestPos, {selected: true}, false)
	    	}
	    }
        function refreshMyLocation(){
        	if(distributor.CurrentPos != null){
        		var pos = distributor.CurrentPos
    	    	setMyLocationMark(pos.Lng, pos.Lat)
        	}	    	
        }
        function flagOrderNodeMarker(){
        	_.each(distributor.AcceptedOrders, function(order){
        		var pos = order.GeoSrc
        		var point = _.findWhere(mapData.Points, {Lng: pos.Lng, Lat:pos.Lat})
        		if(null != point){
        			addPointMarkerToMap(point, {hasOrder: true, orderSigned: order.Signed}, true)
        		}else{
        			console.warn("没有找到订单所在的点，系统异常")
        		}
        	})
        }
	    function drawRouteNodeOnMap(data){
	    	if(data != null){
	    		_.each(data.Points, function(point){ addPointMarkerToMap(point, {type: point.PointType}, true) })
	    		_.each(data.Lines, addLineOverlayToMap)
	    	}
	    }
	    function resetPie(){
        	if(_.size(distributor.AcceptedOrders) <= 0){
        		pie()
        	}else{
        		//对接收的订单按照区域进行分类
        		var groups = _.groupBy(distributor.AcceptedOrders,function(order){return order.Region.Color})
        		var values = _.map(groups,function(v,key){return {value: _.size(v), color: "rgb("+key+")"}})
        		pie(values)
        	}	    	
	    }
	    function addLineOverlayToMap(line){
			var start = line.Start
			var end = line.End
			var points = [new BMap.Point(start.Lng, start.Lat), new BMap.Point(end.Lng, end.Lat)];
			var lineOverlay = new BMap.Polyline(points, {strokeColor:"#50E3C2", strokeWeight:1, strokeOpacity:0.5}); //创建弧线对象
			map.addOverlay(lineOverlay); //添加到地图中
			line.lineOverlay = lineOverlay	    	
	    }
	    function addPointMarkerToMap(point, opt, bBindToPoint){
			var kind = _.findWhere(iconKinds, opt)
			// var kind = _.findWhere(iconKinds, {type: point.PointType})
			if(kind == null){
				console.warn("没找到合适的初始化marker的信息")
				return
			}
			var imageUrl = "/images/marker/"+ kind.imageName
			var myIcon = new BMap.Icon(imageUrl, new BMap.Size(kind.width, kind.height), kind.opt);
			var bmapPoint = new BMap.Point(point.Lng, point.Lat);
		    var marker = new BMap.Marker(bmapPoint, {icon: myIcon});  //创建标注
		    if(null != point.marker){
		    	map.removeOverlay(point.marker)
		    }
		    if(bBindToPoint == true){
			    point.marker = marker
		    }
		    map.addOverlay(marker);                 // 将标注添加到地图中	    
		    return marker	
	    }
	    function setMyLocationMark(lng, lat){
	    	if(myLocationMarker == null){    		
		        var imageUrl = "/images/marker/mylocation.gif"
		        // var myIcon = new BMap.Icon(imageUrl, new BMap.Size(52, 52), {anchor: new BMap.Size(12, 12)});
		        var myIcon = new BMap.Icon(imageUrl, new BMap.Size(64, 64), {anchor: new BMap.Size(16, 16), imageSize: new BMap.Size(32,32)});
		        // var bmapPoint = new BMap.Point(116.644691, 39.934758);//北京物资学院
		        var bmapPoint = new BMap.Point(lng, lat);
		        myLocationMarker = new BMap.Marker(bmapPoint, {icon: myIcon});  //创建标注
		        map.addOverlay(myLocationMarker);                 // 将标注添加到地图中
		        myLocationMarker.setTop(true)
	    	}else{
	    		myLocationMarker.setPosition(new BMap.Point(lng, lat))
	    	}
	    }

	    function refreshNodeToSelect(){
	        //查找可以走向的路径节点，目标点不计算在内
	        //这里有两种情况，正处于路径节点上和在两个节点之间
	        //对于第一种情况，应该查找所有与该点相关的路径
	        //对于第二种情况，显示所在路径的起点与终点
	        if(isDistributorOnNode(distributor) == true){
	        	console.info("配送员在节点上")
	            var pos = distributor.CurrentPos
	            //查找与当前点相关的路线
	            var lines = _.filter(mapData.Lines, function(l){
	                return  _.some([l.Start, l.End], function(p){
	                    return p.Lat == pos.Lat && p.Lng == pos.Lng
	                })
	            })
	            console.log("filter lines :", lines)

	            var points = _.map(lines, function(l){
	                return [l.Start, l.End]
	            })
	            points = _.chain(points).flatten().filter(function(p){
	                return p.Lat != pos.Lat || p.Lng != pos.Lng//去除当前节点
	            }).value()
	            if(distributor.DestPos != null){//去除目标点
	            	var dest = distributor.DestPos
	            	points = _.filter(points, function(p){
	            		return p.Lat != dest.Lat || p.Lng != dest.Lng
	            	})
	            }
	            console.log("points selection: ", points)
	            // addAimMarkers(points)
	            nextNodeSelector = new pointsAimAtLooper(points)
	            nextNodeSelector.next()
	        }else{
	        	console.info("配送员在路上")
	            console.log("起点：",distributor.StartPos)
	            console.log("终点：",distributor.DestPos)
	            nextNodeSelector = new pointsAimAtLooper([distributor.StartPos])
	            nextNodeSelector.next()
	            // addAimMarkers([distributor.StartPos, distributor.DestPos])
	        }          
	    }
	    function onSignOrder(){
	    	//如果有订单，提醒签收
	    	var pos = distributor.CurrentPos
	    	var orderTemp = _.find(distributor.AcceptedOrders, function(order){
	    		return order.GeoSrc.Lat == pos.Lat && order.GeoSrc.Lng == pos.Lng && order.Signed == false
	    	})
	    	if(orderTemp != null){
	    		console.log("签收订单 %s", orderTemp.ID)
		        send({MessageType: {{.pro_sign_order_request}}, Data: {OrderID: orderTemp.ID, DistributorID: distributor.ID}})//请求重置目标点
	    	}
	    }
        function onFocusOnNextNode(){
    		nextNodeSelector.next()
        }
        function onFocusOnPreNode(){
        	nextNodeSelector.pre()
        }
	    function isDistributorOnNode(){
	        var crt = distributor.CurrentPos
    		return _.some([distributor.DestPos, distributor.StartPos], function(point){
    			return point != null && crt.Lat == point.Lat && crt.Lng == point.Lng
    		})
	    }
	    function selectOrder(){
	        var index = mySwiper.activeIndex
	        if(index >= 0){
	            console.log("选择了第 %d 个Slide", index)
	            var $slide = $(mySwiper.slides[index])
	            var $title = $(".slide-title", $slide)
	            var orderID = $title.text()
	            console.log("获取的订单ID为：%s", orderID)
	            var msg = {MessageType: {{.pro_order_select_response}}, Data:{OrderID: orderID, DistributorID: distributorID}}
	            send(msg)

	            // mySwiper.removeSlide(index)
	            // mySwiper.appendSlide('<div class="swiper-slide"> <span class="slide-title">订单编号04</span> <span class="slide-content">地址04</span> </div>')
	        }
	    }
	    function onChooseDestNode(){
	    	var selectedPos = nextNodeSelector.currentPostion()
	    	if(selectedPos == null){
	    		console.info("没有选择到点")
	    	}else{
	    		console.info("选择了一个点作为目标点 ", selectedPos)
	    		var p = _.findWhere(mapData.Points, {Lat: selectedPos.Lat, Lng: selectedPos.Lng})
	    		if(p == null){
	    		    console.warn("没有查找到选中的点")
	    		    return
	    		}else{
	    		    send({MessageType: {{.pro_reset_destination_request}}, Data: {PositionID: p.ID, DistributorID: distributor.ID}})//请求重置目标点
	    		}
	    		// myApp.alert('操作完成', '');
	    		// setTimeout(function () {
	    		//     myApp.closeModal()
	    		// }, 800);

	    		// myApp.showPreloader('已选择')
    		 //    setTimeout(function () {
    		 //        myApp.hidePreloader();
    		 //        myApp.closeModal()
    		 //    }, 800);
	    	}
	    }

	    function onPreparedToStartGame(){
        	viewRouteToPage(mainView, 'waiting')
	        	// viewRouteToPage(mainView, 'process1')
        	send({MessageType: {{.pro_prepared_for_select_order}}, Data:{DistributorID: distributorID}})
	    }
	    function resetOrderSignButtonState(state){
	    	if(state){
	    		$$("#btnSignOrder").removeClass("disabled")
	    	}else{
	    		$$("#btnSignOrder").addClass("disabled")
	    	}
	    }
	    function pie(pieData){
	        if(pieData == null){
	            pieData = [
	                        {
	                            value: 1,
	                            color:"#C8C8C8",
	                            highlight: "#C8C8C8",
	                            label: "无"
	                        }
	                    ];
	        }
	        var ctx = document.getElementById("chart-area").getContext("2d");
	        var myPie = new Chart(ctx).Pie(pieData);    
	    }
	    function send(msg){
	        if (!conn) {
	            return false;
	        }
	        conn.send(JSON.stringify(msg))
	    }
	    function transformTimeElapseToStandardFormat(i){
	    	return String.format("{0}:{1}",padLeft(Math.floor(i/60).toFixed(0), 2), padLeft((i%60).toFixed(0), 2))
	    }
	    function padLeft(str, lenght) {
            if (str.length >= lenght)
                return str;
            else
                return padLeft("0" + str, lenght);
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