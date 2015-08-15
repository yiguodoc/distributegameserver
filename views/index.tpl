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
  <body>
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
		        
		        <div data-page="process1" class="page" id="">
		              <!-- Scrollable page content -->
		              <div class="page-content "> 
				        <div class="content-block" style="margin-top: 20px;  margin-bottom: 15px;">
			                <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
			                <div id="canvas-holder" style="text-align: center;">
	                			<canvas id="chart-area" width="150" height="150"/>
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
				                     <a href="#process1" class="button button-big button-fill" id="" onclick="selectOrder()">选择订单</a>
		                		</div>
		                		<div class="col-10"></div>
		                    </div>
		                </div>

		              </div>
		        </div>

		        <div data-page="index" class="page" id="">
		              <!-- Scrollable page content -->
		              <div class="page-content "> 
				        <div class="content-block" style="margin-top: 100px;">
				                <!-- <p style="text-align: center;">我是{{.distributor.Name}}</p> -->
				                <p style="text-align: center;">我准备好了</p>
				                <div class=" login-btn-content">
				                      <a href="#process1" class="button button-big button-fill" id="" onclick="viewRouteToPage(mainView, 'process1')">进入游戏</a>
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
				        <div id="allmap" style="height:200px;margin-top:10px;"></div>

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
		          <span class="tabbar-label">道具</span>
		      </a>
		  </div>
		</div> 


    </div>
	<script src="javascripts/jquery.js"></script>
	<script type="text/javascript" src="http://api.map.baidu.com/api?type=quick&ak=kU4NWwyP5SwguC2W2WAfO1bO&v=1.0"></script>
    <!-- Path to Framework7 Library JS-->
    <script type="text/javascript" src="framework7/dist/js/framework7.min.js"></script>
    <script type="text/javascript" src="javascripts/lodash.js"></script>
    <script type="text/javascript" src="javascripts/Chart.min.js"></script>
    <!-- Path to your app js-->
    <script type="text/javascript" src="javascripts/app.js"></script>
    <script type="text/javascript">
	    // Initialize your app
	    var myApp = new Framework7();

	    // Export selectors engine
	    var $$ = Dom7;

	    var marker = null
	    var map = new BMap.Map("allmap");
	    var conn;
	    var distributorID = "{{.distributor.ID}}"
	    var orders = []

	    // Add view
	    var mainView = myApp.addView('.view-main', {
	        // Because we use fixed-through navbar we can enable dynamic navbar
	        // dynamicNavbar: true
	    });
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

        var pieData = [
        			{
        				value: 1,
        				color:"#F7464A",
        				highlight: "#FF5A5E",
        				label: "Red"
        			}
        			

        		];

			var ctx = document.getElementById("chart-area").getContext("2d");
			var myPie = new Chart(ctx).Pie(pieData);

		var mySwiper = myApp.swiper('.swiper-container', {
		  pagination: '.swiper-pagination',
		  paginationHide: false,
		  paginationClickable: true,
		  nextButton: '.swiper-button-next',
		  prevButton: '.swiper-button-prev',
		}); 

		function selectOrder(){
			var index = mySwiper.activeIndex
			if(index >= 0){
				var slide = mySwiper.slides[index]			
		        console.log("选择了第 %d 个Slide", index)
		        console.log(slide)
		        var $slide = $(slide)
		        var $title = $(".slide-title", $slide)
		        console.log("获取的订单ID为：%s", $title.text())
		        mySwiper.removeSlide(index)
		        mySwiper.appendSlide('<div class="swiper-slide"> <span class="slide-title">订单编号04</span> <span class="slide-content">地址04</span> </div>')


			}
		}

    </script>
  </body>
</html>    