const host = 'http://' + location.host;
// const host = 'http://127.0.0.1:8888';

Date.prototype.Format = function(fmt) {
	var o = {
		"M+": this.getMonth()+1,                 //月份
		"d+": this.getDate(),                    //日
		"h+": this.getHours(),                   //小时
		"m+": this.getMinutes(),                 //分
		"s+": this.getSeconds(),                 //秒
		"q+": Math.floor((this.getMonth()+3)/3), //季度
		"S": this.getMilliseconds()             //毫秒
	};
	if (/(y+)/.test(fmt)) {
		fmt = fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length));
	}
	for(var k in o) {
		if(new RegExp("("+ k +")").test(fmt)) {
			fmt = fmt.replace(RegExp.$1, (RegExp.$1.length===1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));
		}
	}
	return fmt;
}

$(document).ready(() => {
	const pop = $('#pop');
	const game = $('#game');
	
	$('#play').on('click', () => {
		pop.fadeIn();
		game.attr('src', 'play.html');
	});
	
	pop.on('click', () => {
		pop.fadeOut();
		setTimeout(() => {
			game.attr('src', '');
		}, 300);
	});
	
	const loading = $('#loading');
	const init = $('#init');

	$.ajax({
		type: 'GET',
		url: host + '/api/init',
		success: result => {
			loading.hide();
			const { code, data, note } = result;
			if (code === 0) {
				if (data.init) {
					if (data.remote) {
						alert('远端无法进行初始化，请在本地进行操作');
					} else {
						init.css('display', 'flex');
						$('#name').text(data.name);
						$('#secret').text(data.secret);
						if (data.qrcode) {
							$('#enter').text('使用谷歌验证器APP扫描上方二维码');
							const qrcode = $('#qrcode');
							qrcode.attr('src', data.qrcode);
							qrcode.show();
						}
					}
				} else check();
			} else alert(note);
		},
		complete: () => loading.hide(),
	});
	
	$('#confirm').on('click', () => {
		if (confirm('确认扫描完成吗？二维码将不会再次展示')) {
			loading.show();
			$.ajax({
				type: 'POST',
				url: host + '/api/init',
				success: result => {
					const { code, note } = result;
					if (code === 0) {
						init.hide();
						check();
					} else alert(note);
				},
				complete: () => loading.hide(),
			});
		}
	});
	
	const finish = $('#finish');
	const form = $('#form');
	const forwards = $('#forwards');
	const time = $('#time');

	let srvTime = 0;

	const check = () => {
		loading.show();
		$.ajax({
			type: 'GET',
			url: host + '/api/check',
			success: result => {
				const { code, data, note } = result;
				if (code === 0) {
					if (data.exist) {
						finish.show();
						if (data.forwards.length) {
							for (let forward of data.forwards) {
								forwards.append('<li>' + forward.replace(':', '->') + '</li>')
							}
						} else forwards.hide();
					} else {
						srvTime = data.time;
						refreshTime();
						form.css('display', 'flex');
					}
				} else alert(note);
			},
			complete: () => loading.hide(),
		});
	};

	setInterval(() => {
		srvTime += 1000;
		refreshTime();
	}, 1000)

	// 更新时间显示
	const refreshTime = () => {
		time.text((new Date(srvTime)).Format("yyyy-MM-dd hh:mm:ss"));
	}

	const code = $('#code');
	
	$('#validate').on('click', () => {
		codeVal = code.val();
		if (/^\d{6}$/.test(codeVal)) {
			loading.show();
			$.ajax({
				type: 'POST',
				url: host + '/api/validate',
				data: JSON.stringify({ code: parseInt(codeVal) }),
				success: result => {
					const { code, data, note } = result;
					if (code === 0) {
						form.hide();
						finish.show();
						if (data.length) {
							for (let forward of data) {
								forwards.append('<li>' + forward.replace(':', '->') + '</li>')
							}
						} else forwards.hide();
					} else alert(note);
				},
				complete: () => loading.hide(),
			});
		} else alert('请输入六位校验码');
	});
});
