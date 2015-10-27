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
  <body class="theme-lightblue"  distributorID = "{{.distributor.ID}}">
    <!-- Status bar overlay for full screen mode (PhoneGap) -->
    <div class="statusbar-overlay"></div>
    <!-- Panels overlay-->
    <div class="panel-overlay"></div>
    <!-- Left panel, let it be with reveal effect -->
    <div class="panel panel-left panel-reveal">
        <!-- 左侧栏 -->
        <div class="content-block">
        	<div id = "distributorName" style="  color: rgba(200,200,200,1); font-size: 20px; border-bottom: 1px solid rgba(200,200,200,0.7); padding-bottom: 10px;">{{.distributor.Name}}</div>
              <!-- Click on link with "close-panel" class will close panel -->
              <p><a href="#" class="">历史排名</a></p>
              <p><a href="#" class="">设置</a></p>
              <p><a href="#" class="">关于</a></p>
              <!-- <p><a href="#" class="close-panel">关于</a></p> -->
              <!-- Click on link with "open-panel" and data-panel="right" attribute will open Right panel -->
              <!-- <p><a href="#" data-panel="right" class="open-panel">Open Right Panel</a></p> -->
        </div>
    </div>
    <!-- Views -->
    <!-- <div class="views"> -->
    <div class="views tabs toolbar-through">

		<!-- Your main view, should have "view-main" class -->
		<!-- 信息页面 -->
		<div class="view view-main tab active " id="view-main">
			<!-- <div class="navbar">
			    <div class="navbar-inner" id= "statusNavbar">
			        <div class="center sliding">配送大师</div>
			        <div class="right"> <a href="#" class="link icon-only open-panel"> <i class="icon icon-bars"></i> </a> </div>
			    </div>
			</div> -->
		      <!-- Pages container, because we use fixed-through navbar and toolbar, it has additional appropriate classes-->
		    <div class="pages navbar-fixed">
		        <!-- Page, "data-page" contains page name -->
		        <!-- 订单选择页面 -->
		        <div data-page="processSelectOrder" class="page no-swipeback navbar-fixed" id="1">
		        	<div class="navbar">
		        	    <div class="navbar-inner" >
		        	    	<div class="left" > <a href="#" class="link icon-only back"> <i class="icon icon-back"></i> </a> </div>
		        	        <div class="center sliding">配送大师</div>
		        	        <div class="right"> <a href="#" class="link icon-only open-panel"> <i class="icon icon-bars"></i> </a> </div>
		        	    </div>
		        	</div>		
		            <div class="page-content "> 

				        <div class="content-block" style="margin-top: 20px;  margin-bottom: 15px;">
			                <!-- <div id="canvas-holder" style="text-align: center;">
	                			<canvas id="chart-area" width="130" height="130"/>
	                		</div> -->
			                <!-- <p style="text-align: center;margin-top:0px;font-size: 12px;">订单区域分布比例</p> -->

	                		<div class="card">
                				<div style="text-align: center; padding-top: 10px; font-size: 14px; border-bottom: solid 1px rgba(100,100,100,0.2); padding-bottom: 8px;"> 订单统计</div>
		                		<div class="card-content">
		                		    <div class="card-content-inner" style="padding-top: 0px;">
		                		    	<div class="content-block-inner" style="margin-top:0px;">
		                		    		<div class="row no-gutter" style="margin-top:5px;">
		                		    			<div class="col-50" style="text-align: left;padding-left: 32px;"> 已选订单 </div>
		                		    			<div id="ordersSelectedCount" class="col-50" style="text-align: left;padding-left: 30px;"> 0</div>
		                		    		</div>
		                		    		<div class="row no-gutter" style="margin-top:5px;">
		                		    			<div class="col-50" style="text-align: left;padding-left: 32px;"> 订单分值 </div>
		                		    			<div id = "ordersTotalScore" class="col-50" style="text-align: left;padding-left: 30px;"> 0 </div>
		                		    		</div>
		                		    		
		                		    		<div class="row no-gutter" style="margin-top:5px;">
		                		    			<div class="col-50" style="text-align: left;padding-left: 32px;"> 平均距离 </div>
		                		    			<div id = "ordersAverageDistance" class="col-50" style="text-align: left;padding-left: 30px;"> 0 </div>
		                		    		</div>

		                		    	</div>

		                		    </div>
  	                		    </div>
 	                		</div>  

			            </div>
			            <div class="swiper-custom" style="border-top: 1px solid rgba(100,100,100,0.1);border-bottom: 1px solid rgba(100,100,100,0.1);">
			               <div class="swiper-container">
			                 <div class="swiper-pagination"></div>
			                 <div class="swiper-wrapper">
								<!-- <div class="swiper-slide">
				                   	<span class="slide-title">订单编号01</span>
				                   	<span class="slide-content">地址01</span>
				                   	<span class="slide-content">地址01</span>
			                   </div> -->

			                    
			                 </div>
			               </div>
			               <div class="swiper-button-prev"></div>
			               <div class="swiper-button-next"></div>
			             </div>
		                <div class=" login-btn-content">
						<!-- 		                	<div class="row" style="margin-left: 20px; margin-right: 20px;">
		                		<div class="col-50">
				                     <a href="#" class="button button-raised " id="" onclick="onStartDistribution()">返回配送状态</a>
		                		</div>

		                		<div class="col-50">
				                     <a href="#" class="button button-raised " id="btnSelectOrder" onclick="selectOrder()">选择当前订单</a>
		                		</div>
		                    </div> -->


		                	 <div class="row">
		                		<div class="col-10"></div>
		                		<div class="col-80">
				                     <a href="#" class="button button-big button-fill" id="btnSelectOrder" onclick="selectOrder()">选择订单</a>
		                		</div>
		                		<div class="col-10"></div>
		                    </div>

		                	<!--<div class="row">
		                		<div class="col-20"></div>
		                		<div class="col-60">
				                     <a href="#" class="button " id="" onclick="selectOrder()">开始配送</a>
		                		</div>
		                		<div class="col-20"></div>
		                    </div> -->


		                </div>

		            </div>
		        </div>
		        <!-- 等待其他参与者进入的页面 -->
                <div data-page="waiting" class="page  no-swipeback" id="2">
                	<div class="navbar">
                	    <div class="navbar-inner">
                	        <div class="center sliding">配送大师</div>
                	    </div>
                	</div>	
                      <!-- Scrollable page content -->
                      <div class="page-content "> 
        		        <div class="content-block" style="margin-top: 100px;">
        		                <div id="waitingInfoList">
	        		                <!-- <p id="waitingInfo" style="text-align: center;">等待其他人进入...</p> -->
	        		                <!-- <p id="waitingInfo" style="text-align: center;">等待其他人进入...</p> -->
        		                </div>
        		                <div class="content-block-title" style="font-size: 13px; color: rgba(100,100,100,0.5); text-align: center;">等待配送员进入</div>
    		                    <div class="list-block">
        		                    <ul id="waitingListBox">
        		                      <!-- <li class="item-content" > <div class="item-inner" style="text-align: center;"> <div class="item-title" style="width:100%;">Item title</div> </div> </li> -->

        		                    </ul>
    		                    </div>
        		                <!-- <div class=" login-btn-content">
        		                      <a href="#processSelectOrder" class="button button-big button-fill" id="" onclick="viewRouteToPage(mainView, 'processSelectOrder')">进入游戏</a>
        		                </div> -->

        	            </div>

                      </div>
                </div>
                <!-- 订单选择完毕，转入配送状态页面之前的页面 -->
                <div data-page="processGo2Distribution" class="page  no-swipeback" id="4">
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
                <div data-page="processDistribution" class="page  no-swipeback" id="5"><!-- Scrollable page content -->
                    <div class="navbar">
                        <div class="navbar-inner" >
                            <div class="center sliding">配送大师</div>
                            <div class="right"> <a href="#" class="link icon-only open-panel"> <i class="icon icon-bars"></i> </a> </div>
                        </div>
                    </div>	   
					<div class="page-content "> 
					    <div class="content-block" style="margin-top: 0px;">
					            <!-- <p style="text-align: center;">00:10</p> -->
					            <div id="systimeShow" class="" style="text-align: center;margin-top:5px;  color: rgba(150,150,150,0.9);"></div>
					            <div class="content-block-inner" style="text-align: center;margin-top:50px;">
					            	<div style="font-size: 13px; color: rgba(150,150,150,0.8);">订单签收进度</div>
					            	<div id="orderSignProcess" style="font-size: 19px;">7/10</div>

					            </div>
					            <div style = "  text-align: center; margin-bottom: 10px; font-size: 28px; border-bottom: 1px solid rgba(200,200,200,0.5);">
					            	<!-- <span>配送中</span> -->
					            	<img id="imgStateGif" src = "/images/marker/running.gif" style="width:100px;margin-top: 20px;margin-bottom: -8px;">
					            </div>
					            <div style="text-align: center;margin-bottom: 20px; color: rgba(100,100,100,0.7);"> <span id="speed">0km/h</span> </div>
					            <div class=" login-btn-content">
					                  <!-- <a href="#" class="button button-big button-fill disabled" id="btnSignOrder" onclick="onSignOrder()">订单签收</a> -->

                                  	<div class="row" style="margin-left: 20px; margin-right: 20px;">
                                  		<div class="col-50">
                  		                     <a href="#" class="button button-raised disabled" id="btnSignOrder" onclick="onSignOrder()">订单签收</a>
                                  		</div>
                                  		<div class="col-50">
                  		                     <a href="#" class="button button-raised " onclick="onRouteToSelectOrderView()">订单选择</a>
                                  		</div>

                                    </div>
					            </div>
					            <div style = "text-align:center;margin-top: 25px; text-decoration: underline;"> <span onclick="onEndDistribution()">结束配送 </span> </div>
					    </div>

					</div>
                </div>
                <!-- 配送完成统计页面 -->
                <div data-page="processStatistic" class="page  no-swipeback" id="6"><!-- Scrollable page content -->
                       
					<div class="page-content "> 
					    <div class="content-block" style="margin-top: 0px;">
				            <!-- <p style="text-align: center;">00:10</p> -->
				            <div class="" style="text-align: center; margin-top: 55px; color: rgba(100,100,100,1.9); font-size: 20px;">订单配送已完成</div>
				            <div style = "text-align: center; font-size: 15px; border-bottom: 1px solid rgba(200,200,200,0.4); padding-top: 40px; color: rgba(100,100,100,0.5);">
				            	<span>成绩统计</span>
				            </div>
				            <div class="content-block-inner" style="margin-top:0px;">
				            	<div class="row no-gutter" style="margin-top:5px;">
				            		<div class="col-50" style="text-align: left;padding-left: 70px;"> 排名 </div>
				            		<div id="statisticRank" class="col-50" style="text-align: left;padding-left: 50px;"> 0</div>
				            	</div>
				            	<div class="row no-gutter" style="margin-top:5px;">
				            		<div class="col-50" style="text-align: left;padding-left: 70px;"> 总得分 </div>
				            		<div id = "statisticScoreTotal" class="col-50" style="text-align: left;padding-left: 50px;"> 0 </div>
				            	</div>
				            	
				            	<div class="row no-gutter" style="margin-top:5px;">
				            		<div class="col-50" style="text-align: left;padding-left: 70px;"> 总耗时 </div>
				            		<div id = "statisticTimeTotal" class="col-50" style="text-align: left;padding-left: 50px;"> 0 </div>
				            	</div>

				            </div>
					    </div>
					</div>
                </div>

                <!-- 登录进入游戏之后的页面 -->
		        <div data-page="index" class="page  no-swipeback" id="3">
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



				        <div id="allmap" style="height:99%;margin-top:1px;width:100%"> </div>
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
				        <div  class="content-block-title" style="margin-top: 20px;">未签收<span id="orderListUnsignedLabel">（0）</span></div>
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

		<!-- tab工具栏 -->
		<div class="toolbar tabbar tabbar-labels">
		  <div class="toolbar-inner">
		      <a href="#view-main" class="tab-link active">
		          <i class="icon tabbar-icon-status"></i>
		          <span class="tabbar-label">状态</span>
		      </a>
		      <a href="#view-map" class="tab-link" onclick="onMapActive()">
		          <i class="icon tabbar-icon-map">
		              <!-- <span class="badge bg-red">5</span> -->
		          </i>
		          <span class="tabbar-label">地图</span>
		      </a>
		      <a href="#view-orders" class="tab-link">
		          <i class="icon tabbar-icon-order"></i>
		          <span class="tabbar-label">订单</span>
		      </a>
		      <a href="#view-cards" class="tab-link">
		          <i class="icon tabbar-icon-msg"></i>
		          <span class="tabbar-label">消息</span>
		      </a>
		  </div>
		</div> 


    </div>
    <script type="text/javascript" src="http://api.map.baidu.com/api?v=2.0&ak=kU4NWwyP5SwguC2W2WAfO1bO"></script>

    <script type="text/javascript">
    	var restartingGame = false
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
	    	{MessageType: {{.pro_2c_game_start}}, handler: pro_2c_game_start_handler},
	    	{MessageType: {{.pro_2c_message_broadcast_before_game_start}}, handler: pro_2c_message_broadcast_before_game_start_handler},
	    	{MessageType: {{.pro_2c_order_distribution_proposal}}, handler: pro_2c_order_distribution_proposal_handler},
	    	{MessageType: {{.pro_timer_count_down}}, handler: pro_timer_count_down_handler},
	    	{MessageType: {{.pro_2c_message_broadcast}}, handler: pro_2c_message_broadcast_handler},
	    	{MessageType: {{.pro_2c_order_select_result}}, handler: pro_2c_order_select_result_handler},
	    	{MessageType: {{.pro_2c_reach_route_node}}, handler: pro_2c_reach_route_node_handler},
	    	{MessageType: {{.pro_2c_move_to_new_position}}, handler: pro_2c_move_to_new_position_handler},
	    	{MessageType: {{.pro_2c_move_from_node}}, handler: pro_2c_move_from_node_handler},
	    	{MessageType: {{.pro_2c_map_data}}, handler: pro_2c_map_data_handler},
	    	{MessageType: {{.pro_2c_distributor_info}}, handler: pro_2c_distributor_info_handler},
	    	{MessageType: {{.pro_2c_reset_destination}}, handler: pro_2c_reset_destination_handler},
	    	{MessageType: {{.pro_2c_sign_order}}, handler: pro_2c_sign_order_handler},
	    	{MessageType: {{.pro_2c_all_order_signed}}, handler: pro_2c_all_order_signed_handler},
	    	{MessageType: {{.pro_2c_speed_change}}, handler: pro_2c_speed_change_handler},
	    	{MessageType: {{.pro_2c_end_game}}, handler: pro_2c_end_game_handler},
	    	{MessageType: {{.pro_2c_rank_change}}, handler: pro_2c_rank_change_handler},
	    	{MessageType: {{.pro_2c_restart_game}}, handler: pro_2c_restart_game_handler},
	    	{MessageType: {{.pro_2c_on_line_user_change}}, handler: pro_2c_on_line_user_change_handler},
	    	{MessageType: {{.pro_2c_sys_time_elapse}}, handler: pro_2c_sys_time_elapse_handler, print: false},
	    	{}
	    ]
	    function pro_2c_on_line_user_change_handler(msg){
	    	var distributorsOffline = msg.Data
	    	console.info("配送员上线情况发生变化")
	    	var waitingListBox = $$("#waitingListBox")
	    	waitingListBox.children().each(function(index, item){
	    		item.remove()
	    	})
	    	var dom = '<li class="item-content" > <div class="item-inner" style="text-align: center;"> <div class="item-title" style="width:100%;">{0}</div> </div> </li>'
	    	_.each(distributorsOffline, function(d){
		    	waitingListBox.append(String.format(dom, d.Name))
		    	console.log(d.Name + " 不在线")
	    	})
	    }
	    function pro_2c_game_start_handler(msg){
	    	distributor = msg.Data
	    	refreshDistributionStateView()//数据
        	viewRouteToPage(mainView, 'processDistribution')

	    }
	    function pro_2c_restart_game_handler(msg){
            // window.location.href = window.location.href
            // window.location.reload();
            restartingGame = true
	    }
	    function pro_2c_rank_change_handler(msg){
	    	distributor = msg.Data
	    	refreshStatisticPage()
	    }
	    function pro_2c_end_game_handler(msg){
			distributor = msg.Data
			refreshStatisticPage()
			//提示游戏结束
			myApp.alert('游戏结束，将转入统计页面', '提示');
        	viewRouteToPage(mainView, 'processStatistic')	    	
	    }
	    //速度发生变化
	    function pro_2c_speed_change_handler(msg){
			distributor = msg.Data
	    	// refreshSpeed()
	    	refreshDistributionStateView()
	    }
    	//更新系统时间
	    function pro_2c_sys_time_elapse_handler(msg){
	    	var time = transformTimeElapseToStandardFormat(msg.Data)
	    	$$("#systimeShow").text(time)
	    	// console.log("系统时间更新：", time)
	    }
	    function pro_2c_all_order_signed_handler(msg){
			distributor = msg.Data
			myApp.modal({
			  title:  '提示',
			  text: '订单签收完毕，是否结束游戏查看成绩？',
			  buttons: [
			    {
			      text: '是',
			      onClick: function() {
			        send({MessageType: {{.pro_end_game_request}}, Data: null})//请求重置目标点
			      }
			    },
			    {
			      text: '否',
			      onClick: function() {
			        
			      }
			    },
			  ]
			})
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
				refreshOrderView()
				// refreshOrderSignProcess()
				refreshDistributionStateView()
	    	}
	    }
	    function pro_2c_reset_destination_handler(msg){
	    	console.info("重置目的地成功")
    		distributor = msg.Data
	    	setDestinationMarker()
			// refreshNodeToSelect()
	    }
	    function pro_2c_reach_route_node_handler(msg){
    		distributor = msg.Data
			refreshMyLocation()
			console.info("配送员到达节点")
			refreshNodeToSelect()
	    	setDestinationMarker()
	    	remindOrderSigning()
	    	// refreshRunningState(2)
	    	refreshDistributionStateView()
	    }

	    function pro_2c_move_to_new_position_handler(msg){
    		// distributor = msg.Data
    		//Position
    		distributor.CurrentPos = msg.Data
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
	    	// refreshRunningState(1)
	    	refreshDistributionStateView()
	    }
	    function pro_2c_message_broadcast_before_game_start_handler(msg){
	    	var html = '<p id="waitingInfo" style="text-align: center;">'+msg.Data+'</p>'
	    	$$("#waitingInfoList").prepend(html)
	    }
	    function pro_2c_order_distribution_proposal_handler(msg){
	    	mySwiper.removeAllSlides();
	    	var orders = msg.Data
	    	var currentPostion = distributor.CurrentPos

	    	_.each(orders, function(order){
	    		var pos = order.GeoSrc
	    		var distance = CoolWPDistance(currentPostion.Lat, currentPostion.Lng, pos.Lat, pos.Lng).toFixed(0) + "米"
	    		var dom =   '<div class="swiper-slide"> '+
					    		'<span class="slide-title" style="margin-left: 10px; margin-right: 10px;">{0}</span> '+
					    		'<span class="slide-content" style="margin-top:20px;font-size: 13px;">{1}</span> '+
					    		'<div class="row no-gutter" style="margin-top:5px;"> <div class="col-50" style="text-align: right;padding-right: 10px;font-size: 13px;"> 分值 </div> <div  class="col-50" style="text-align: left;padding-left: 10px;font-size: 13px;"> {2} </div> </div>'+
					    		'<div class="row no-gutter" style="margin-top:5px;"> <div class="col-50" style="text-align: right;padding-right: 10px;font-size: 13px;"> 距离 </div> <div  class="col-50" style="text-align: left;padding-left: 10px;font-size: 13px;"> {3} </div> </div>'+
				    		'</div>'
	    	    mySwiper.appendSlide(String.format(dom, order.ID, order.GeoSrc.Address+"", order.Score, distance))
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
	    		// resetPie()
	    		refreshSeletedOrdersStatistics()
				setAcceptedOrderMarkerOnMap()
				refreshOrderView()
				var lastOrder = distributor.AcceptedOrders[_.size(distributor.AcceptedOrders)-1]
				addMsgToView({timeStamp: transformTimeElapseToStandardFormat(lastOrder.SelectedTime), content: "选择配送订单 "+ lastOrder.ID})
				// refreshOrderSignProcess()
				refreshDistributionStateView()
	    	}else{
	    	    console.log("没有抢到订单 ")	    		
	    	}
	    }
	    //初始化配送员数据页面
	    function pro_2c_distributor_info_handler(msg){
	    	var data = msg.Data
	    	console.log(data)
	    	distributor = data
			refreshMyLocation()
			setAcceptedOrderMarkerOnMap()
			refreshNodeToSelect()
			remindOrderSigning()
			refreshOrderView()
			addMsgToView({timeStamp: transformTimeElapseToStandardFormat(distributor.TimeElapse), content: "上线"})
			// refreshOrderSignProcess()
			// refreshRunningState()
			// refreshSpeed()
			refreshDistributionStateView()


	    	switch(distributor.CheckPoint){
	    		case {{.checkpoint_flag_origin}}:
		        	viewRouteToPage(mainView, 'index')
	    		break
	    		case {{.checkpoint_flag_prepared_for_game}}:
		        	viewRouteToPage(mainView, 'waiting')//倒计时页面

	    		case {{.checkpoint_flag_game_started}} :
		        	// viewRouteToPage(mainView, 'processSelectOrder')
		        	// refreshSeletedOrdersStatistics()
		        	viewRouteToPage(mainView, 'processDistribution')		        	
	    		break
	    		case {{.checkpoint_flag_game_over}}:
		        	viewRouteToPage(mainView, 'processStatistic')
		        	refreshStatisticPage()
	    		break
	    	}

	    }
	    function onMapActive(){
	    	console.log("onMapActive...")


	    	// if(distributor.CurrentPos != null){
	    	// 	var pos = distributor.CurrentPos
	    	// 	map.setCenter(new BMap.Point(pos.Lng, pos.Lat))
	    	//     // map.setCenter(new BMap.Point(116.404212, 39.914888))
	    	//     // setMapMarker(pos.Lng, pos.Lat, false)
	    	// }
	    	setTimeout(function(){
		    	if(myLocationMarker != null){
		    	    var pos = myLocationMarker.getPosition()
		    	    map.setCenter(pos)
		    	}	    		
	    	}, 300)
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
	    //配送结果统计页面刷新
	    function refreshStatisticPage(){
	    	$$("#statisticRank").text(distributor.Rank)
	    	$$("#statisticScoreTotal").text(distributor.Score)
	    	$$("#statisticTimeTotal").text(transformTimeElapseToStandardFormat(distributor.TimeElapse))
	    }
	    function onEndDistribution(){
	    	var text = "配送结束将统计你的最终成绩，确定吗？"
	    	var orderSigned = _.filter(distributor.AcceptedOrders, function(order){return order.Signed == true})
	    	if(_.size(distributor.AcceptedOrders) > _.size(orderSigned)){
	    		text = "您尚有未签收的订单，结束配送后将不能签收，确定吗？"
	    	}
	    	myApp.modal({
	    	  title:  '提示',
	    	  text: text,
	    	  buttons: [
	    	    {
	    	      text: '是',
	    	      onClick: function() {
	    	        send({MessageType: {{.pro_end_game_request}}, Data: null})//请求重置目标点
	    	      }
	    	    },
	    	    {
	    	      text: '否',
	    	      onClick: function() {
	    	        
	    	      }
	    	    },
	    	  ]
	    	})
	    }
		//刷新配送状态页面
	    function refreshDistributionStateView(){
	    	var speed = distributor.CurrentSpeed
	    	$$("#speed").text(speed+"km/h")	
	    	if(isDistributorOnNode(distributor) == true){
	    		$$("#imgStateGif").attr("src", "/images/marker/stayRunning.gif")//原地跑动
	    	}else{
	    		$$("#imgStateGif").attr("src", "/images/marker/running.gif")//跑动
	    	}	 
	    	var orderSigned = _.filter(distributor.AcceptedOrders, function(order){return order.Signed == true})
	    	$$("#orderSignProcess").text(String.format("{0}/{1}", _.size(orderSigned), _.size(distributor.AcceptedOrders)))	    	   	    	
	    }

	    //添加消息到消息页面
	    function addMsgToView(msg){
	    	var dom = $$("#msgList")
	    	dom.prepend(String.format('<li> <div class="item-inner"> <div class="item-title-row"> <div class="item-title"></div> <div class="item-after" style="font-size:14px;">{0}</div> </div>  <div class="item-text" style="padding-left: 5px; color: black; ">{1}</div> </div> </li>', msg.timeStamp, msg.content))
	    }
	    //重置订单页面的订单列表，当订单数量或者状态发生变化时调用
	    function refreshOrderView(){
	    	var refresher = function(domID, orders){
	    		var dom = $$("#"+domID)
		    	dom.children().remove()

		    	$$("#"+domID+"Label").text(String.format("（{0}）", _.size(orders)))
	    		_.each(orders, function(order){
	    			dom.append(String.format('<li> <a href="#" class="item-link item-content"> <div class="item-inner"> <div class="item-title-row"> <div class="item-title">{0}</div> <div class="item-after"><span class="badge  bg-lightblue" style="padding-top:3px; border-radius: 5px;">+{1}</span></div> </div> <div class="item-text">{2}</div> </div> </a> </li>', order.ID,order.Score, order.Address || ""))
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
        //将配送员选择的订单标识在地图上
        function setAcceptedOrderMarkerOnMap(){
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
	    function refreshSeletedOrdersStatistics(){
	    	var orderCount = _.size(distributor.AcceptedOrders)
	    	$$("#ordersSelectedCount").text(orderCount)
	    	$$("#ordersTotalScore").text(_.reduce(distributor.AcceptedOrders, function(total, order){
	    		return total + order.Score
	    	}, 0))
	    	var currentPostion = distributor.CurrentPos
	    	var totalDistance = _.reduce(distributor.AcceptedOrders, function(total, order){
	    		var pos = order.GeoSrc
	    		var distance = CoolWPDistance(currentPostion.Lat, currentPostion.Lng, pos.Lat, pos.Lng)
	    		return total + distance
	    	}, 0)
	    	if(orderCount > 0){
		    	$$("#ordersAverageDistance").text((totalDistance/orderCount).toFixed(0)+"米")
	    	}else{
		    	$$("#ordersAverageDistance").text(0)

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
	    function isDistributorOnNode(){
	        var crt = distributor.CurrentPos
    		return _.some([distributor.DestPos, distributor.StartPos], function(point){
    			return point != null && crt.Lat == point.Lat && crt.Lng == point.Lng
    		})
	    }
	    function selectOrder(){
	    	var debunced = _.debounce(function(){
		        var index = mySwiper.activeIndex
		        if(index >= 0){
		            console.log("选择了第 %d 个Slide", index)
		            var $slide = $(mySwiper.slides[index])
		            var $title = $(".slide-title", $slide)
		            var orderID = $title.text()
		            console.log("获取的订单ID为：%s", orderID)
		            var msg = {MessageType: {{.pro_order_select_response}}, Data:{OrderID: orderID, DistributorID: distributorID}}
		            send(msg)
		        }	    		
	    	}, 500)
	    	debunced()
	    	$$("#btnSelectOrder").addClass("disabled")
	    	setTimeout(function(){
	    		$$("#btnSelectOrder").removeClass("disabled")
	    	}, 5000)
	    }
	    function onRouteToSelectOrderView(){
        	viewRouteToPage(mainView, 'processSelectOrder')
	    	// $$("#statusNavbar").prepend('<div id="statusNavbarLeft" class="left"> <a href="#" class="link icon-only"> <i class="icon icon-back"></i> </a> </div>')
        	$$("#statusNavbarLeft").show()
	    	
	    }

	    function onStartDistribution(){

	    	mainView.router.back()
        	// viewRouteToPage(mainView, 'processDistribution')
        	$$("#statusNavbarLeft").hide()
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