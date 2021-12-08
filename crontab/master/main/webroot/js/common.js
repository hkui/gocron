/**
 * Created by huangkui@lepu.cn on 2019/8/12.
 */

$.ajaxSettings.timeout = 5000;
$.ajaxSettings.xhrFields = {
    withCredentials: true
};
$.ajaxSettings.error = function (event, status, settings) {
    if (status == 'timeout') {
        layer.open({content: '请求超时，请稍后再试', time: 1});
    } else {
        layer.open({content: '请求异常，请稍后再试', time: 1});
    }
};
$.ajaxSettings.complete = function (event, status) {
    if (status == 'success') {
        process(event.responseJSON)
    }

}

function process(res) {
    if (res.errno == 2) {
        setTimeout(function () {
            window.location.href = '/login.html'
        }, 1500)

    }
}
