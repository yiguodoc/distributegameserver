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
  <title>创建新游戏</title>
</head>

<body style="height:100%;">
  {{template "topNav.tpl" .}}

  <div class="sui-layout" style="height:85%;margin-top:0px;">
    {{template "sidebar.tpl" .}}
    <div class="content" style="height:100%;margin-left: 195px; margin-right: 5px;border-left: 3px solid rgba(100,100,100,0.3); padding-left: 2px;">
      <div style="  border-top: solid 1px rgba(100,100,100,0.3);  border-bottom: solid 1px rgba(100,100,100,0.3); height:80px;">
        <div style="font-size: 18px; margin-top: 10px; padding-left: 5px;">创建一个新游戏</div>
        <div style="font-size: 13px; margin-top: 20px; padding-left: 5px;color: rgba(150,150,150,0.8);">选择一张地图和列表中的成员，创建一个新游戏</div>
      </div>
      <a href="javascript:void(0);" onclick="startGame()" class="sui-btn btn-xlarge btn-success" style="width:150px;  margin-top: 25px; margin-left: 8px;">创建</a>
      <div>
        <!-- <div>选择地图</div> -->
        <div style="margin-bottom: 30px;padding-left: 8px; margin-top: 30px;">
          <div style="font-size: 14px; margin-bottom: 5px; color: rgba(100,100,100,0.6);">选择地图</div>
          <div style="width:80%; border-bottom: dotted 1px rgba(100,100,100,0.2);margin-bottom: 10px; margin-top: 5px;"></div>
          <select id="selectMap" style="width: 150px; height: 22px;">
            <!-- <option value="saab">Saab</option> -->
          </select>
          <div style="width:40%; border-bottom: dotted 1px rgba(100,100,100,0.2);margin-bottom: 5px; margin-top: 15px;color: rgba(100,100,100,0.6);">地图信息</div>
          <div style="margin-top: 0px;">
            <span style="color:gray;">游戏时长(分钟)：</span>
            <span id="mapInfoGameLength"></span>
            <span style="margin-left:60px;color:gray;">订单数：</span>
            <span id="mapInfoOrderCount"></span>
            <span style="margin-left:60px;color:gray;">出生点：</span>
            <span id="mapInfoBornPointCount">
          </div>
        </div>
        <div style="margin-bottom: 30px;padding-left: 8px; margin-top: 30px;">
          <div style="font-size: 14px; margin-bottom: 5px; color: rgba(100,100,100,0.6);">选择模式</div>
          <div style="width:50%; border-bottom: dotted 1px rgba(100,100,100,0.2);margin-bottom: 10px; margin-top: 5px;"></div>
          <select id="selectMode" style="width: 150px; height: 22px;">
            <option value="dual">个人对抗</option>
            <option value="team">团队任务</option>
          </select>
        </div>

        <div style="font-size: 14px; margin-bottom: 5px; color: rgba(100,100,100,0.6); padding-left: 8px;">选择成员</div>
        <div style="width:100%; border-bottom: dotted 1px rgba(100,100,100,0.2);margin-bottom: 10px; margin-top: 5px; margin-left:8px;"></div>
        <div style="margin-bottom:-20px; margin-top: -25px;">
          <a href="javascript:void(0);" onclick="refreshUserTable()" class="sui-btn btn-bordered btn-info" style="width:150px;  margin-top: 25px; margin-left: 8px;">刷新</a></br>
        </div>
        <!-- <div style="border-bottom: solid 1px rgba(0,0,0,0.1); margin-top: 5px;"></div> -->
        <table id="dtProcess" class="display" cellspacing="0" width="100%">
          <thead>
            <tr>
              <th>ID</th>
              <th>姓名</th>
              <th>团队</th>
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

      $.get("/mapNameList", function(data) {
        console.log("map list => ", data.Data)
        var mapNameList = data.Data
        var $selectMapList = $("#selectMap")
        _.each(mapNameList, function(name) {
          $selectMapList.append('<option value="' + name + '">' + name + '</option>')
        })
        refreshMapDetail()
      })
      $("#selectMap").bind("change", refreshMapDetail)

      table = $('#dtProcess').DataTable({
        "columnDefs": [],
        "paging": false,
        "ordering": true,
        "order": [
          [0, "asc"]
        ],
        "info": false,
        "searching": true,
        "ajax": {
          "url": "/distributors?atgame=0",
          "dataSrc": "Data"
        },
        "columns": [{
          "data": "UserInfo.ID",
          "width": "20%"
        }, {
          "data": "UserInfo.Name",
          "width": "40%"
        }, {
          "data": "UserInfo.Team",
          "width": "40%"
        }, ]
      });

      $('#dtProcess tbody').on('click', 'tr', function() {
        //只能选中单行
        if ($(this).hasClass('selected')) {
          $(this).removeClass('selected');
        } else {
          // table.$('tr.selected').removeClass('selected');
          $(this).addClass('selected');
        }

        //可以选中多行
        // $(this).toggleClass('selected');
      });
    })

    function getSelectedRows() {
      var data = table.rows(".selected").data()
      console.log(data)
      return data
    }

    function startGame() {
      var rows = getSelectedRows()
      if (rows.length <= 0) {
        $.alert("先选择一些成员参与该游戏吧")
        return
      }
      console.log("selected to be in game :")
      _.each(rows, function(row) {
        console.log(row.ID + " " + row.Name)
      })
      var idList = _.reduce(rows, function(list, row) {
          list.push(row.UserInfo.ID)
          return list
        }, [])
        // var postData = {id: JSON.stringify(idList)}//form

      var $options = $("#selectMap option:selected")
      if ($options.length <= 0) {
        alert("需要先选择地图")
        return
      }

      var selectedOption = $options[0]
      var mode = $("#selectMode option:selected").val()
      console.log("select map => " + selectedOption.value)
      console.log("select mode => " + mode)
      var postData = {
        id: idList,
        mapID: selectedOption.value,
        mode: mode
      }
      console.log(postData)
      postData = JSON.stringify(postData)
      $.ajax({
          url: "/newGame",
          type: "POST",
          dataType: "json",
          contentType: "application/json; charset=utf-8",
          data: postData,
          success: function(data) {
            console.log(data)
            if (data != null) {
              if (data.Code == 0) {
                alert("游戏创建成功，可以在游戏列表中查看！")
                refreshUserTable()
              } else {
                alert("创建失败：" + data.Message)
              }
            }
          }
        })
        // $.post("/newGame", postData, function(data) {
        //     console.log(data)
        //     if (data != null) {
        //         if (data.Code == 0) {

      //         }
      //     }
      // }, "json")

    }

    function refreshUserTable() {
      table.ajax.reload()
    }

    function refreshMapDetail() {
      //下载地图信息
      var selectedMapID = $("#selectMap  option:selected").val()
      if (selectedMapID != null && selectedMapID.length > 0) {
        $.get("/mapData?id=" + selectedMapID, function(data) {
          if (data.Code == 0) {
            var mapData = data.Data
            $("#mapInfoGameLength").text(transformTimeElapseToStandardFormat(mapData.TimeLength))
            $("#mapInfoOrderCount").text(_.size(_.where(mapData.Points, {
              PointType: 2
            })))
            $("#mapInfoBornPointCount").text(_.size(_.where(mapData.Points, {
              IsBornPoint: true
            })))
          }
        })
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
