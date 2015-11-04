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
    
    .statusLabel {
        margin-top: 20px;
        padding-left: 5px;
        font-size: 14px;
        color: rgba(100, 100, 100, 0.8);
        /*border-bottom: 1px solid rgba(100,100,100,0.2);*/
        padding-bottom: 3px;
    }
    </style>
    <script src="javascripts/jquery.js"></script>
    <script type="text/javascript" src="http://g.alicdn.com/sj/dpl/1.5.1/js/sui.min.js"></script>
    <script src="/dataTable/js/jquery.dataTables.js"></script>
    <script src="javascripts/lodash.js"></script>
    <script src="javascripts/string-format.js"></script>
    <title>比分排名</title>
</head>

<body style="height:100%;">
    {{template "topNav.tpl" .}}

    <div class="sui-layout" style="height:85%;margin-top:0px;">
        {{template "sidebar.tpl" .}}

        <div class="content" style="height:100%;margin-left: 195px; margin-right: 5px;border-left: 3px solid rgba(100,100,100,0.3); padding-left: 2px;">
            <div style="  border-top: solid 1px rgba(100,100,100,0.3);  border-bottom: solid 1px rgba(100,100,100,0.3); height:80px;">
                <div style="font-size: 18px; margin-top: 10px; padding-left: 5px;">详细信息</div>
                <div style="font-size: 13px; margin-top: 20px; padding-left: 5px;color: rgba(150,150,150,0.8);">查看游戏的状态，地图和参与者的排名信息</div>
            </div>
                <a href="javascript:void(0);" onclick="refresh()" class="sui-btn btn-xlarge btn-success" style="width:100px; margin-top: 16px; margin-left: 5px;">刷新</a></br>

            <div style="margin-top: 40px;">
                <div style="" class="statusLabel"> 
                    <span>当前状态</span> 
                    <span id = "gameStatus" style="color: black; font-size: 20px; padding-left: 10px;">进行中</span>
                    <span id = "gameProcess" style="padding-left: 15px;">(2:30/5:00)</span>
                </div>
                <div style="" class="statusLabel"> 
                    <span>地图</span> 
                    <span id = "mapName" style="color: black; font-size: 20px; padding-left: 40px;">繁忙都市</span> 
                </div>
            </div>
            <!-- <div id="currentState">繁忙都市</div> -->
            <!-- <div id = "currentState">进行中 2:30/5:00</div> -->
            <div style="" class="statusLabel">排名信息</div>
            <div>
                <div style="border-bottom: solid 1px rgba(0,0,0,0.1); margin-top: 10px;"></div>
                <table id="dtProcess" class="display" cellspacing="0" width="100%">
                    <thead>
                        <tr>
                            <th>排名</th>
                            <th>ID</th>
                            <th>姓名</th>
                            <th>得分</th>
                            <th>耗时</th>
                        </tr>
                    </thead>
                </table>
            </div>
        </div>
    </div>
    <div style="width: 100%; text-align: center; font-size: 13px; padding-top: 10px; border-top: 1px solid rgba(100,100,100,0.3); color: rgba(100,100,100,0.8); margin-top: 5px;">配送大师团队技术支持</div>
    <script type="text/javascript">
    var gameID = "{{.gameID}}"
    var table
    var gameStatusList = [{key: 0, value: "尚未开始"}, {key: 1, value: "进行中"}, {key: 2, value: "已结束"}]

    $(function() {
        refresh()
    })

    function refresh() {
        $.get("/gameList?gameID={{.gameID}}", function(msg){
            if(msg.Data == null || _.size(msg.Data) <= 0){
                return 
            }
            var unit = msg.Data[0]
            console.log(unit)
            var distributors = unit.Distributors
            if(table != null){
                table.destroy()
            }
            $("#gameStatus").text(_.result(_.findWhere(gameStatusList, {key: unit.Status}), 'value'))
            $("#gameProcess").text(unit.TimeElapse  + "/" + unit.BasicInfo.Game_time_loop * 60)
            $("#mapName").text(unit.BasicInfo.MapName)
            table = $('#dtProcess').DataTable({
                "data": distributors,
                "columnDefs": [{
                        "render": function(data, type, row) {
                            return transformTimeElapseToStandardFormat(row.TimeElapse)
                        },
                        "targets": 4
                    },

                ],
                "paging": false,
                "ordering": true,
                "order": [
                    [0, "asc"]
                ],
                "info": false,
                "searching": true,
                "columns": [{
                    "data": "Rank",
                    "width": "20%"
                }, {
                    "data": "UserInfo.ID",
                    "width": "20%"
                }, {
                    "data": "UserInfo.Name",
                    "width": "20%"
                }, {
                    "data": "Score",
                    "width": "20%"
                }, {
                    // "data": "TimeElapse",
                    "width": "20%"
                }, ]
            });
        })


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
