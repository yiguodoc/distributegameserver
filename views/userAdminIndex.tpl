<!DOCTYPE html>
<html  style="height:100%;">
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
    <title>学生管理</title>
</head>


<body style="height:100%;">
    {{template "topNav.tpl" .}}

    <div class="sui-layout" style="height:100%;margin-top:0px;">
        {{template "sidebar.tpl" .}}
        <div class="content" style="height:100%;margin-left: 195px; margin-right: 5px;border-left: 3px solid rgba(100,100,100,0.3); padding-left: 2px;">
            <div style="  border-top: solid 1px rgba(100,100,100,0.3);  border-bottom: solid 1px rgba(100,100,100,0.3); height:80px;">
                <div style="font-size: 18px; margin-top: 10px; padding-left: 5px;">学生列表</div>
                <div style="font-size: 13px; margin-top: 20px; padding-left: 5px;color: rgba(150,150,150,0.8);">管理你的学生</div>
            </div>
            <div>
                <div style="margin-bottom:-18px;">
                    <a href="javascript:void(0);" onclick="refresh_grid()" class="sui-btn btn-xlarge btn-bordered btn-info" style="width:100px;  margin-top: 25px; margin-left: 8px;">刷新</a>
                    <a href="javascript:void(0);" onclick="add()" class="sui-btn btn-xlarge btn-success" style="width:100px;  margin-top: 25px; margin-left: 8px;">增加</a>
                    <a href="javascript:void(0);" onclick="resetpwd()" class="sui-btn btn-xlarge btn-success" style="width:100px;  margin-top: 25px; margin-left: 8px;">重置密码</a>
                    <a href="javascript:void(0);" onclick="deleteUser()" class="sui-btn btn-xlarge btn-success" style="width:100px;  margin-top: 25px; margin-left: 8px;">删除</a> 
                    <a href="javascript:void(0);" onclick="groupUser()" class="sui-btn btn-xlarge btn-success" style="width:100px;  margin-top: 25px; margin-left: 8px;">组队</a>
                    <a href="javascript:void(0);" onclick="leaveGroup()" class="sui-btn btn-xlarge btn-success" style="width:100px;  margin-top: 25px; margin-left: 8px;">离开团队</a>
                </div>
                <!-- <div style="border-bottom: solid 1px rgba(0,0,0,0.1); margin-top: 5px;"></div> -->
                <table id="dtProcess" class="display" cellspacing="0" width="100%">
                    <thead>
                        <tr>
                            <th>学号</th>
                            <th>姓名</th>
                            <th>团队</th>
                        </tr>
                    </thead>
                </table>
            </div>
        </div>
    </div>

 
    <script type="text/javascript">
    var table = null
    $(function() {
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
                "url": "/users",
                "dataSrc": "Data"
            },
            "columns": [{
                "data": "ID",
                "width": "20%"
            }, {
                "data": "Name",
                "width": "40%"
            }, {
                "data": "Team",
                "width": "40%"
            }]
        });
        $('#dtProcess tbody').on('click', 'tr', function() {
            //只能选中单行
            if ($(this).hasClass('selected')) {
                $(this).removeClass('selected');
            } else {
                // table.$('tr.selected').removeClass('selected');
                $(this).addClass('selected');
            }

        });

    })

    //open modal
    function add() {
        $("#myModal").modal("show")
    }

    function refresh_grid() {
        table.ajax.reload()
    }

    function getSelectedID() {
        var data = table.rows(".selected").data()
        if (data.length <= 0) {
            alert("请先选择一个用户！")
            return []
        } else {
            return _.pluck(data, "ID")
            // var id = data[0].ID
            // return id
        }
    }

    function deleteUser() {
        var list = getSelectedID()
        if(_.size(list) > 0) {
            console.log("delete users ", list)
            var body = '删除该信息将无法恢复，确定删除吗？'
            if(_.size(list) > 1){
                body = '删除该信息将无法恢复，而且是同时删除多条信息，确定删除吗？'
            }
            $.confirm({
                body: body,
                width: 'normal',
                backdrop: true,
                bgcolor: 'none',
                okHide: function() {
                    $.ajax({
                        url: '/users',
                        type: 'DELETE',
                        contentType: "application/json;charset=UTF-8",
                        data: JSON.stringify(list),
                        success: function(data) {
                            if (data.Code != 0) {
                                $.alert(data.Message)
                            } else {
                                refresh_grid()
                            }
                        }
                    });
                },
                cancelHide: function() {
                    console.log('cancelHide')
                },
            })
        }
    }

    function resetpwd() {
        var list = getSelectedID()
        console.log("resetpwd users ", list)
        if (_.size(list) > 0) {
            var body = '该用户密码将会被重置为默认密码，请尽快修改该密码以保证安全！'
            if(_.size(list) > 1){
                body = '多个用户的密码将会被重置为默认密码，请尽快修改该密码以保证安全！'
            }
            $.confirm({
                body: body,
                width: 'normal',
                backdrop: true,
                bgcolor: 'none',
                okHide: function() {
                    $.ajax({
                        url: "/resetpwd", 
                        type: "PATCH",
                        contentType: "application/json;charset=UTF-8",
                        processData: false,
                        data: JSON.stringify(list),
                        success: function(data) {
                            if (data.Code != 0) {
                                alert(data.Message)
                            }else{
                                alert('重置密码成功！')
                            }
                        }
                    })
                },
            })

        }
    }
    </script>

    <!-- Modal-->
    <div id="myModal" tabindex="-1" role="dialog" data-hasfoot="false" class="sui-modal hide fade">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" data-dismiss="modal" aria-hidden="true" class="sui-close">×</button>
                    <h4 id="myModalLabel" class="modal-title">添加新用户</h4>
                </div>
                <div class="modal-body">
                    <form class="sui-form form-horizontal">
                        <div class="control-group">
                            <label for="inputEmail" class="control-label">学号：</label>
                            <div class="controls">
                                <input type="text" id="inputID" placeholder="1001">
                            </div>
                        </div>
                        <div class="control-group">
                            <label for="inputEmail" class="control-label">姓名：</label>
                            <div class="controls">
                                <input type="text" id="inputName" placeholder="名称">
                            </div>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" data-ok="modal" class="sui-btn btn-primary btn-large" onclick="addUser()">添加</button>
                    <button type="button" data-dismiss="modal" class="sui-btn btn-default btn-large" onclick="cancel()">取消</button>
                </div>
            </div>
        </div>
        <script type="text/javascript">
        function cancel() {
            $("#myModal").modal("hide")
        }

        function addUser() {
            var obj = {
                id: $("#inputID").val(),
                name: $("#inputName").val()
            }
            $.post("/users", obj, function(data) {
                if (data.Code == 0) {
                    refresh_grid()
                } else {
                    alert(data.Message)
                }
                cancel()
            })
        }
        </script>
    </div>
</body>

</html>
