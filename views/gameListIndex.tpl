<!DOCTYPE html>
<html lang="en" style="height:100%;">

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no" />
    <link rel="icon" type="image/png" href="/images/logo.png">
    <link href="http://g.alicdn.com/sj/dpl/1.5.1/css/sui.min.css" rel="stylesheet">
    <link href="/dataTable/css/jquery.dataTables.css" rel="stylesheet" media="screen">
    <style type="text/css">
    td.highlight {
        background-color: rgba(0, 256, 0, 0.1) !important;
    }
    
    tr.highlight {
        background-color: rgba(0, 256, 0, 0.1) !important;
    }
    
    #dtProcess {
        text-align: center;
    }
    
    th {
        text-align: center;
    }
    
    .title {
        font-size: 32px;
        font-weight: 800;
        margin-bottom: 10px;
    }
    
    .itemActive {
        background-color: rgba(0, 0, 0, 0.1);
    }
    
    .itemUnactive {
        background-color: rgba(0, 0, 0, 0);
    }
    </style>
    <script src="javascripts/jquery.js"></script>
    <script type="text/javascript" src="http://g.alicdn.com/sj/dpl/1.5.1/js/sui.min.js"></script>
    <script src="/dataTable/js/jquery.dataTables.js"></script>
    <script src="javascripts/lodash.js"></script>
    <script src="javascripts/string-format.js"></script>
    <title>游戏列表</title>
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
                <!-- <div style="  margin-bottom: 20px; margin-top: 10px; font-size: 16px; border-bottom: solid 1px rgba(100,100,100,0.2); padding-bottom: 10px; color: rgba(100,100,100,0.7);">新游戏</div> -->
                <a href="/newGameIndex" class="sui-btn btn-xlarge btn-info" style="width:150px;  margin-top: 25px; margin-left: 8px;">新游戏</a></br>
                <a href="/gameListIndex" class="sui-btn btn-xlarge btn-info" style="width:150px;  margin-top: 25px; margin-left: 8px;">游戏列表</a></br>
            </div>
        </div>
        <div class="content" style="height:100%;margin-left: 195px; margin-right: 5px;border-left: 3px solid rgba(100,100,100,0.3); padding-left: 2px;">
            <div style="  border-top: solid 1px rgba(100,100,100,0.3);  border-bottom: solid 1px rgba(100,100,100,0.3); height:80px;">
                <div style="font-size: 18px; margin-top: 10px; padding-left: 5px;">正在进行中的游戏</div>
                <div style="font-size: 13px; margin-top: 20px; padding-left: 5px;color: rgba(150,150,150,0.8);">选择一个游戏，查看其进度详情</div>
            </div>
            <div>
                <div style="margin-bottom:-18px;">
                    <a href="javascript:void(0);" onclick="GameDetail()" class="sui-btn btn-xlarge btn-success" style="width:150px;  margin-top: 25px; margin-left: 8px;">详情</a>
                    <a href="javascript:void(0);" onclick="refreshData()" class="sui-btn btn-xlarge btn-bordered btn-info" style="width:150px;  margin-top: 25px; margin-left: 8px;">刷新</a></br>
                </div>
                <!-- <div style="border-bottom: solid 1px rgba(0,0,0,0.1); margin-top: 5px;"></div> -->
                <table id="dtProcess" class="display" cellspacing="0" width="100%">
                    <thead>
                        <tr>
                            <th>游戏ID</th>
                            <th>游戏时长</th>
                            <th>参与人数</th>
                            <th>地图名称</th>
                        </tr>
                    </thead>
                </table>
            </div>
        </div>
    </div>
    <div style="width: 100%; text-align: center; font-size: 13px; padding-top: 10px; border-top: 1px solid rgba(100,100,100,0.3); color: rgba(100,100,100,0.8); margin-top: 5px;">配送大师团队技术支持</div>
    <script type="text/javascript">
    var table = null;
    $(function() {
        table = $('#dtProcess').DataTable({
            "columnDefs": [
                {
                    "render": function ( data, type, row ) {
                        return transformTimeElapseToStandardFormat(row.GameTimeMaxLength)
                    },
                    "targets": 1
                },
                {
                    "render": function ( data, type, row ) {
                        return _.size(row.Distributors)
                    },
                    "targets": 2
                },
            ],
            "paging": false,
            "ordering": true,
            "order": [
                [0, "asc"]
            ],
            "info": false,
            "searching": true,
            "ajax": {
                "url": "/gameList",
                "dataSrc": "Data"
            },
            "columns": [{
                "data": "ID",
                "width": "20%"
            }, {
                // "data": "Name",
                "width": "20%"
            }, {
                // "data": "Name",
                "width": "20%"
            }, {
                "data": "MapName",
                "width": "20%"
            }, ]
        });

        $('#dtProcess tbody').on('click', 'tr', function() {
            //只能选中单行
            if ($(this).hasClass('selected')) {
                $(this).removeClass('selected');
            } else {
                table.$('tr.selected').removeClass('selected');
                $(this).addClass('selected');
            }

            //可以选中多行
            // $(this).toggleClass('selected');
        });
    })
    function refreshData(){
        table.ajax.reload()
    }
    function getSelectedRows() {
        var data = table.rows(".selected").data()
        console.log(data)
        return data
    }
    function GameDetail(){
        var rows = getSelectedRows()
        if(rows.length > 0){
            window.location.href = "/rankIndex?gameID="+rows[0].ID
        }
    }


    function transformTimeElapseToStandardFormat(i) {
        return String.format("{0}:{1}", padLeft(Math.floor(i / 60).toFixed(0), 2), padLeft((i % 60).toFixed(0), 2))
    }

    function padLeft(str, lenght) {
        if (str.length >= lenght)
            return str;
        else
            return padLeft("0" + str, lenght);
    }

    String.format = function() {
        if (arguments.length == 0) {
            return null;
        }
        var str = arguments[0];
        for (var i = 1; i < arguments.length; i++) {
            var re = new RegExp('\\{' + (i - 1) + '\\}', 'gm');
            str = str.replace(re, arguments[i]);
        }
        return str;
    }

    function appendLog(message) {
        console.log(message)
    }
    </script>
</body>

</html>
