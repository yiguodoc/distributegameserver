<!DOCTYPE html>

<html >
<head>
<title>配送员列表</title>
<script src="javascripts/jquery.js"></script>
<script src="javascripts/underscore.js"></script>
<script language="javascript" type="text/javascript">

    var output;
    $(function() {
        // output = document.getElementById("output");
        output = $("#output")
        showUsersList()

    })
    function showUsersList(){
        $.get("/distributors",function(data){
            console.log(data)
            subscribers = data
            // subscribers = JSON.parse(data)
            // console.log(subscribers)
            _.each(subscribers, function(sub){
                var ele = '<a href="/orderDistribute?id='+sub.ID+'">'+sub.Name+'</a></br>'
                // writeToScreen(ele)
                output.append(ele);
                // $("#output").append(ele)
            })
        })
    }

    function init() {
        output = document.getElementById("output");
        showUsersList()
    }
    function writeToScreen(message) {
        var pre = document.createElement("p");
        pre.style.wordWrap = "break-word";
        pre.innerHTML = message;
        output.appendChild(pre);
    }
    // window.addEventListener("load", init, false);
</script>
</head>
<body>
    <h2>配送员列表（测试）</h2>
    <div id="output">
    </div>
</body>
</html>
  
  