<!DOCTYPE html>
<html lang="en">
<head>
<title>比分排名</title>
<script src="javascripts/jquery.js"></script>
<script src="javascripts/lodash.js"></script>
<script src="javascripts/string-format.js"></script>
<style type="text/css">
        body, html,#allmap {width: 100%;height: 100%;overflow: hidden;margin:0;font-family:"微软雅黑";}

</style>
</head>
<body>
    <div> <a href="#" onclick="restartGame()">重新开始 </a> </div>
    <div> <a href="#" onclick="refreshRankData()">刷新排名 </a> </div>
    <div id="output"></div>

</body>
<script type="text/javascript">
    var output;
    $(function() {
        output = $("#output")


    })
    function restartGame(){
        $.get("/restartGame", function(data){

        })
    }
    function refreshRankData(){
        output.empty()
        $.get("/distributors", function(data){
            console.log(data)
            var distributors = _.sortBy(data, "Rank")
            _.each(distributors, function(d){
                // appendLog(d.Name + " : " + d.Rank)
                appendLog(format("{Name} : 排名：{Rank}  得分：{Score}  用时：{TimeElapse}", d))
            })
        })
    }
    function appendLog(message) {
        output.prepend(message+"</br>");
    }    
</script>


</html>
