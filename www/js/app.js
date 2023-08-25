var app = ajsf('sysinfo', (context, rootElement) => {
	context.sysinfo = {
		cpu: {
			cores: []
		},
		ram: {
			fgColor: 'green'
		},
		swap: {
			fgColor: 'red'
		},
		tempDevices: []
	};

	context.loadData = () => {
		$.get('sysinfo.json', {
			success: data => {
				context.data = JSON.parse(data);

				context.sysinfo.ram.label = context.data.ram.used.toFixed(1) + '/' + context.data.ram.total.toFixed(1) + ' G';
				context.sysinfo.ram.percent = context.data.ram.usage;
				context.sysinfo.ram.legend = "RAM";
				context.sysinfo.swap.label = context.data.swap.used.toFixed(1) + '/' + context.data.swap.total.toFixed(1) + ' G';
				context.sysinfo.swap.percent = context.data.swap.usage;
				context.sysinfo.swap.legend = "SWAP";

				for (var i = 0; i < context.data.cpu.cores.length; i++) {
					context.sysinfo.cpu.cores[i] = {
						identifier: context.data.cpu.cores[i].id,
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
							label: context.data.temps[i].sensors[j].temp.toFixed(0) + "°C",
							legend: context.data.temps[i].sensors[j].name,
							percent: context.data.temps[i].sensors[j].temp/130*100
						};
					}
				}

				context.makeRefresh();
			},
			error: () => {
				console.log('Could not load sysinfo');
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

app.directive('circle',`<div class="circle"><div class="circle-inner" ajsf-text="model.label"></div><span class="circle-legend" ajsf-text="model.legend"></span></div>`, (context, el) => {
	context.onRefresh = () => {
		if (context.model.fgColor == undefined) {
			context.model.fgColor = "#008800";
		}
		if (context.model.bgColor == undefined) {
			context.model.bgColor = "#888888";
		}

		var size = $(el).attr("size");
		if (size == undefined) {
			size = 2
		}

		$(el).find('.circle').attr('style', "--progress: " + context.model.percent +"; --fgColor: " + context.model.fgColor + "; --bgColor: " + context.model.bgColor + "; --size: " + size + "rem;");
	};
});

app.directive('bar', `<div class="bar"><div class="bar-inner" ajsf-text="model.identifier"></div><span class="bar-label" ajsf-text="model.label"></span></div>`, (context, el) => {
	if (context.model.fgColor == undefined) {
		context.model.fgColor = "#008800";
	}
	if (context.model.bgColor == undefined) {
		context.model.bgColor = "#888888";
	}

	$(el).find('.bar').attr('style', "--progress: " + context.model.percent + "%; --fgColor: " + context.model.fgColor + "; --bgColor: " + context.model.bgColor + ";");
});