const host = 'http://' + location.host;
// const host = 'http://127.0.0.1:8888';

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
					init.css('display', 'flex');
					$('#name').text(data.name);
					$('#secret').text(data.secret);
					if (data.qrcode) {
						$('#enter').text('使用谷歌验证器APP扫描上方二维码');
						const qrcode = $('#qrcode');
						qrcode.attr('src', data.qrcode);
						qrcode.show();
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
					} else form.css('display', 'flex');
				} else alert(note);
			},
			complete: () => loading.hide(),
		});
	};
	
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
