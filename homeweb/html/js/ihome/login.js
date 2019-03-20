function getCookie(name) {
    var r = document.cookie.match("\\b" + name + "=([^;]*)\\b");
    return r ? r[1] : undefined;
}

$(document).ready(function () {
    $("#email").focus(function () {
        $("#email-err").hide();
    });
    $("#password").focus(function () {
        $("#password-err").hide();
    });
    $(".form-login").submit(function (e) {
        // 阻止浏览器对于表单的默认提交行为
        e.preventDefault();
        var email = $("#email").val();
        var emailPat = /^(.+)@(.+)$/;
        var matchArray = email.match(emailPat);
        var passwd = $("#password").val();
        if (!email || matchArray == null) {
            $("#email-err span").html("请填写正确的邮箱！");
            $("#email-err").show();
            return;
        }
        if (!passwd) {
            $("#password-err span").html("请填写密码!");
            $("#password-err").show();
            return;
        }
        // 将表单的数据存放到对象data中
        var data = {};
        $(this).serializeArray().map(function (x) {
            data[x.name] = x.value;
        });
        // 将data转为json字符串
        var jsonData = JSON.stringify(data);
        $.ajax({
            url: "/api/v1.0/sessions",
            type: "post",
            data: jsonData,
            contentType: "application/json",
            dataType: "json",
            headers: {
                "X-CSRFTOKEN": getCookie("csrf_token"),
            },
            success: function (data) {
                if ("0" == data.errno) {
                    // 登录成功，跳转到主页
                    location.href = "/";
                    return;
                }
                else {
                    // 其他错误信息，在页面中展示
                    $("#password-err span").html(data.errmsg);
                    $("#password-err").show();
                    return;
                }
            }
        });
    });
})