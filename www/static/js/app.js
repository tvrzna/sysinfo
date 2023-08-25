var app = ajsf('sysinfo', (context, rootElement) => {
	context.sysinfo = {
		cpu: {
			cores: []
		},
		memory: [],
		tempDevices: [],
		diskusage: []
	};

	var lock = false;
	var loadingOverlay = $(rootElement).find('.loading-spinner');

	context.addLoading = (count) => {
		if (count == undefined) {
			count = 1;
		}
		context.loading += count;

		if (context.loading > 0) {
			setTimeout(() => {
				if (context.loading > 0) {
					loadingOverlay[0].style.display = 'flex';
				}
			}, 250)
		} else {
			loadingOverlay[0].style.display = 'none';
			context.loading = 0;
		}
	};

	context.loadData = () => {
		if (lock) {
			return;
		}
		lock = true;
		context.addLoading();
		$.get('sysinfo.json', {
			success: data => {
				context.data = JSON.parse(data);
				context.sysinfo.memory = [
					{
						label: context.data.ram.used.toFixed(1) + '/' + context.data.ram.total.toFixed(1) + ' G',
						percent: context.data.ram.used / context.data.ram.total * 100,
						legend: 'RAM'
					},
					{
						label: context.data.swap.used.toFixed(1) + '/' + context.data.swap.total.toFixed(1) + ' G',
						percent: context.data.swap.used / context.data.swap.total * 100,
						legend: 'SWAP'
					},
				];

				for (var i = 0; i < context.data.cpu.cores.length; i++) {
					context.sysinfo.cpu.cores[i] = {
						legend: context.data.cpu.cores[i].id,
						percent: context.data.cpu.cores[i].usage,
						label: context.data.cpu.cores[i].usage.toFixed(0) + '% ' + context.data.cpu.cores[i].mhz.toFixed(0) + 'MHz'
					};
				}

				context.sysinfo.loadavg1 = context.data.loadavg.loadavg1.toFixed(2);
				context.sysinfo.loadavg5 = context.data.loadavg.loadavg5.toFixed(2);
				context.sysinfo.loadavg15 = context.data.loadavg.loadavg15.toFixed(2);

				context.sysinfo.tempDevices = [];
				for (var i = 0; i < context.data.temps.length; i++) {
					context.sysinfo.tempDevices[i] = {
						name: context.data.temps[i].name,
						temps: []
					};
					for (var j = 0; j < context.data.temps[i].sensors.length; j++) {
						context.sysinfo.tempDevices[i].temps[j] = {
							label: context.data.temps[i].sensors[j].temp.toFixed(0) + '°C',
							legend: context.data.temps[i].sensors[j].name,
							percent: context.data.temps[i].sensors[j].temp/130*100
						};
					}
				}

				context.sysinfo.diskusage = [];
				for (var i = 0; i < context.data.diskusage.length; i++) {
					context.sysinfo.diskusage[i] = {
						label: context.data.diskusage[i].usedgb.toFixed(0) + '/' + context.data.diskusage[i].totalgb.toFixed(0) + ' G',
						legend: context.data.diskusage[i].path,
						percent: context.data.diskusage[i].usedgb/context.data.diskusage[i].totalgb*100
					}
				}

				var days = Math.floor(context.data.uptime / 60 / 60 / 24);
				var hours = Math.floor((context.data.uptime / 60 / 60)) - days * 24;
				var minutes = Math.floor(context.data.uptime / 60) - (days * 24 + hours) * 60;
				var seconds = context.data.uptime - ((days * 24 + hours) * 60 + minutes) * 60;
				context.sysinfo.uptime = (days > 0 ? days + ' days, ' : '') + String(hours).padStart(2, '0') + ':' + String(minutes).padStart(2, '0') + ':' + String(seconds).padStart(2, '0');

				context.makeRefresh();
			},
			error: () => {
				console.log('Could not load sysinfo');
			},
			complete: () => {
				context.addLoading(-1);
				lock = false;
			}
		});
	};

	context.makeRefresh = () => {
		context.refresh();

		var coresHeight = 0;
		$(rootElement).find(".cpu-cores > bar").each(function(i, el){
			coresHeight += el.offsetHeight;
		});

		$(rootElement).find('.cpu-cores').attr("style", "max-height: " + coresHeight/2 + "px;");
	};

	context.loadData();
	setInterval(context.loadData, 1000);
});

app.directive('circle',`<div class="circle" ajsf-title="model.legend | suffix ' '| suffix model.label"><div class="circle-inner" ajsf-text="model.label"></div><span class="circle-legend" ajsf-text="model.legend"></span></div>`, (context, el) => {
	context.onRefresh = () => {
		if (context.model.fgColor == undefined) {
			context.model.fgColor = context.model.percent <= 40 ? "#008800" : (context.model.percent <= 80 ? "#888800" : "#880000");
		}
		if (context.model.bgColor == undefined) {
			context.model.bgColor = "#FFFFFF10";
		}

		var size = $(el).attr("size");
		if (size == undefined) {
			size = 2
		}

		$(el).find('.circle').attr('style', "--progress: " + context.model.percent +"; --fgColor: " + context.model.fgColor + "; --bgColor: " + context.model.bgColor + "; --size: " + size + "rem;");
	};
});

app.directive('bar', `<div class="bar" ajsf-title="model.legend | suffix ' '| suffix model.label"><span class="bar-legend" ajsf-text="model.legend"></span><div class="bar-inner"></div><span class="bar-label" ajsf-text="model.label"></span></div>`, (context, el) => {
	context.onRefresh = () => {
		if (context.model.fgColor == undefined) {
			context.model.fgColor = context.model.percent <= 40 ? "#008800" : (context.model.percent <= 80 ? "#888800" : "#880000");
		}
		if (context.model.bgColor == undefined) {
			context.model.bgColor = "#FFFFFF10";
		}

		var size = $(el).attr("size");
			if (size == undefined) {
				size = 6
			}

		$(el).find('.bar').attr('style', "--progress: " + context.model.percent + "%; --fgColor: " + context.model.fgColor + "; --bgColor: " + context.model.bgColor + "; --size: " + size + "rem;");
	};
});