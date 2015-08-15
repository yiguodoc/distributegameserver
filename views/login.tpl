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
    <!-- <link rel="stylesheet" href="framework7/dist/css/my-app.css"> -->


  </head>
  <body>

    <div class="login-screen">
      <div class="view">
        <!-- Top Navbar-->
        <div class="navbar">
          <div class="navbar-inner">
            <!-- We need cool sliding animation on title element, so we have additional "sliding" class -->
            <div class="center sliding">登录</div>
          </div>
        </div>
        <!-- Pages container, because we use fixed-through navbar and toolbar, it has additional appropriate classes-->
        <div class="pages navbar-through toolbar-through">
              <!-- Page, "data-page" contains page name -->
              <div data-page="index" class="page">
                  <!-- Scrollable page content -->
                  <div class="page-content ">

                        <div class="content-block">
                            <div class="row">
                              <div class="col-10"></div>
                              <div class="col-80"  style="text-align: center;">
                                <img src="/images/logo.png" width="80">
                              </div>
                              <div class="col-10"></div>
                            </div>
                            <div class="logo"  style="text-align: center;">配送大师</div>
                        </div>

                        <div class="content-block login-input-content">
                          <div class="store-data list-block">
                            <ul>
                                <li>
                                  <div class="item-content">
                                    <div class="item-media"><i class="icon icon-form-name"></i></div>
                                    <div class="item-inner">
                                      <!-- <div class="item-title label">Name</div> -->
                                      <div class="item-input">
                                        <input type="text" placeholder="学号">
                                      </div>
                                    </div>
                                  </div>
                                </li>

                                <li>
                                  <div class="item-content">
                                    <div class="item-media"><i class="icon icon-form-password"></i></div>
                                    <div class="item-inner">
                                      <!-- <div class="item-title label">Name</div> -->
                                      <div class="item-input">
                                        <input type="password" placeholder="密码">
                                      </div>
                                    </div>
                                  </div>
                                </li>
                            </ul>
                          </div>
                        </div>

                        <div class="row">
                            <!-- Each "cell" has col-[widht in percents] class -->
                            <div class="col-10"></div>
                            <div class="col-80">
                                <div class=" login-btn-content">
                                      <a href="/index" class="button button-big button-fill external" id="btnLogin">登录</a>
                                </div>

                            </div>
                            <div class="col-10"></div>
                        </div>
                        
                        <div class="row" style="margin-top:15px;">
                            <!-- Each "cell" has col-[widht in percents] class -->
                            <div class="col-25"></div>
                            <div class="col-50" style="text-align: center;">
                                  <a href="javascript:">忘记密码</a>
                            </div>
                            <div class="col-25"></div>
                        </div>
                  </div>
              </div>
          </div>
      </div>
    </div>
    <!-- Status bar overlay for full screen mode (PhoneGap) -->
    <div class="statusbar-overlay"></div>
    <!-- Panels overlay-->
    <div class="panel-overlay"></div>

    <!-- Views -->
    <div class="views">
      <!-- Your main view, should have "view-main" class -->
          <div class="view view-main">
            <!-- Top Navbar-->
            <!--         <div class="navbar">
              <div class="navbar-inner">
                <div class="center sliding"></div>
              </div>
            </div> -->
            <!-- Pages container, because we use fixed-through navbar and toolbar, it has additional appropriate classes-->
            <div class="pages navbar-through toolbar-through">
              <!-- Page, "data-page" contains page name -->
              <div data-page="index" class="page">
                <!-- Scrollable page content -->
                <div class="page-content "> </div>
                    <div class="toolbar">
                        <div class="row">
                          <div class="col-50">
                            <a href="#" class="button button-big open-login-screen ">登录</a>
                          </div>
                          
                          <div class="col-50">
                            <a href="#" class="button button-big">注册</a>
                          </div>
                          
                        </div>
                        <!-- <div class="toolbar-inner"> </div> -->
                    </div>
                </div>
              </div>
          </div>

    </div>
    <!-- Path to Framework7 Library JS-->
    <script type="text/javascript" src="framework7/dist/js/framework7.min.js"></script>
    <!-- Path to your app js-->
    <!-- // <script type="text/javascript" src="javascripts/appLogin.js"></script> -->
    <script type="text/javascript">
    // Initialize your app
    var myApp = new Framework7();

    // Export selectors engine
    var $$ = Dom7;

    // Add view
    var mainView = myApp.addView('.view-main', {
        // Because we use fixed-through navbar we can enable dynamic navbar
        dynamicNavbar: true
    });
    // $$("#btnLogin").on("click", function(e){
    //     console.log("加载首页")
    //     mainView.router.load({url:"/index"})
    // })

    </script>

  </body>
</html>    