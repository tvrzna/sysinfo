var app = ajsf('sysinfo', (context, rootElement) => {
	context.sysinfo = {
		cpu: {
			cores: []
		},
		memory: [],
		tempDevices: [],
		diskusage: []
	};
	context.intervalValue = 1000;

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
			}, 600)
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
				if (context.data.ram != undefined && context.data.swap != undefined) {
					context.sysinfo.memory = [
						{
							label: context.data.ram.used.toFixed(1) + context.data.ram.usedUnit + '/' + context.data.ram.total.toFixed(1) + context.data.ram.totalUnit,
							percent: context.data.ram.percent,
							legend: 'RAM'
						},
						{
							label: context.data.swap.used.toFixed(1) + context.data.swap.usedUnit + '/' + context.data.swap.total.toFixed(1) + context.data.swap.totalUnit,
							percent: context.data.swap.percent,
							legend: 'SWAP'
						},
					];
				}

				if (context.data.cpu != undefined) {
					context.sysinfo.cpu.cores = [];
					for (var i = 0; i < context.data.cpu.cores.length; i++) {
						context.sysinfo.cpu.cores[i] = {
							legend: context.data.cpu.cores[i].id,
							percent: context.data.cpu.cores[i].usage,
							label: context.data.cpu.cores[i].usage.toFixed(0) + '% ' + context.data.cpu.cores[i].mhz.toFixed(0) + 'MHz'
						};
					}
				}

				if (context.data.loadavg != undefined) {
					context.sysinfo.loadavg1 = context.data.loadavg.loadavg1.toFixed(2);
					context.sysinfo.loadavg5 = context.data.loadavg.loadavg5.toFixed(2);
					context.sysinfo.loadavg15 = context.data.loadavg.loadavg15.toFixed(2);
				}

				if (context.data.temps != undefined) {
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
								percent: context.data.temps[i].sensors[j].temp / 130 * 100
							};
						}
					}
				}

				if (context.data.diskusage != undefined) {
					context.sysinfo.diskusage = [];
					for (var i = 0; i < context.data.diskusage.length; i++) {
						context.sysinfo.diskusage[i] = {
							label: context.data.diskusage[i].used.toFixed(0) + context.data.diskusage[i].usedUnit + '/' + context.data.diskusage[i].total.toFixed(0) + context.data.diskusage[i].totalUnit,
							legend: context.data.diskusage[i].path,
							percent: context.data.diskusage[i].percent
						}
					}
				}

				if (context.data.uptime != undefined) {
					var days = Math.floor(context.data.uptime / 60 / 60 / 24);
					var hours = Math.floor((context.data.uptime / 60 / 60)) - days * 24;
					var minutes = Math.floor(context.data.uptime / 60) - (days * 24 + hours) * 60;
					var seconds = context.data.uptime - ((days * 24 + hours) * 60 + minutes) * 60;
					context.sysinfo.uptime = (days > 0 ? days + ' days, ' : '') + String(hours).padStart(2, '0') + ':' + String(minutes).padStart(2, '0') + ':' + String(seconds).padStart(2, '0');
				}

				if (context.data.netspeed != undefined) {
					context.sysinfo.netspeed = [];
					for (var i = 0; i < context.data.netspeed.length; i++) {
						context.sysinfo.netspeed[i] = {
							name: context.data.netspeed[i].name,
							data: [{
								label: context.data.netspeed[i].download.toFixed(1) + context.data.netspeed[i].downloadUnit,
								legend: 'D',
								percent: context.data.netspeed[i].downloadPercent
							},
							{
								label: context.data.netspeed[i].upload.toFixed(1) + context.data.netspeed[i].uploadUnit,
								legend: 'U',
								percent: context.data.netspeed[i].uploadPercent
							}
							]
						}
					}
				}

				if (context.data.top != undefined) {
					context.sysinfo.top = [];
					for (var i = 0; i < context.data.top.length; i++) {
						context.sysinfo.top[i] = {
							pid: context.data.top[i].pid,
							state: context.data.top[i].state,
							comm: context.data.top[i].comm,
							cpu: context.data.top[i].cpu.toFixed(2),
							ram: context.data.top[i].ram.toFixed(2),
							ramUnit: context.data.top[i].ramUnit
						};
					}
				}

				if (context.data.diskstats != undefined) {
					context.sysinfo.diskstats = [];
					for (var i = 0; i < context.data.diskstats.length; i++) {
						context.sysinfo.diskstats[i] = {
							name: context.data.diskstats[i].name,
							iops: [
								{
									label: context.data.diskstats[i].riops,
									legend: 'R iops',
									percent: context.data.diskstats[i].riops * 20 / 100
								},
								{
									label: context.data.diskstats[i].wiops,
									legend: 'W iops',
									percent: context.data.diskstats[i].wiops * 20 / 100
								}
							],
							speed: [
								{
									label: context.data.diskstats[i].read.toFixed(2) + context.data.diskstats[i].readUnit,
									legend: 'R',
									percent: context.data.diskstats[i].readPercent
								},
								{
									label: context.data.diskstats[i].write.toFixed(2) + context.data.diskstats[i].writeUnit,
									legend: 'W',
									percent: context.data.diskstats[i].writePercent
								}
							]
						}
					}
				}

				if (context.data.smartctl != undefined) {
					context.sysinfo.smartctl = [];
					for (var i = 0; i < context.data.smartctl.length; i++) {
						context.sysinfo.smartctl[i] = {
							name: context.data.smartctl[i].name,
							model: context.data.smartctl[i].model,
							smartStatusPassed: context.data.smartctl[i].smartStatusPassed,
							powerOnTime: context.data.smartctl[i].powerOnTime,
							powerCycleCount: context.data.smartctl[i].powerCycleCount,
							temp: context.data.smartctl[i].temp,
							attributes: []
						};

						for (var j = 0; j < context.data.smartctl[i].attributes.length; j++) {
							context.sysinfo.smartctl[i].attributes[j] = {
								name: context.data.smartctl[i].attributes[j].name,
								value: context.data.smartctl[i].attributes[j].value,
								worst: context.data.smartctl[i].attributes[j].worst,
								raw: context.data.smartctl[i].attributes[j].raw,
								flags: context.data.smartctl[i].attributes[j].flags,
							}
						}
					}
				}

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
		var coresCount = 0;
		$(rootElement).find(".cpu-cores > bar").each(function (i, el) {
			coresHeight += el.offsetHeight;
			coresCount++;
		});

		var multiplier = 4;
		if (coresCount >= 16 && coresCount < 32) {
			multiplier = 8;
		} else if (coresCount >= 32 && coresCount < 64) {
			multiplier = 16;
		} else if (coresCount >= 64) {
			multiplier = 32;
		}

		$(rootElement).find('.cpu-cores').attr("style", "max-height: " + coresHeight / Math.floor(coresCount / multiplier) + "px;");
	};

	context.changeInterval = () => {
		if ([1000, 2000, 5000, 10000].includes(Number(context.intervalValue))) {
			clearInterval(context.interval);
			context.interval = setInterval(context.loadData, Number(context.intervalValue));
			console.log('changed interval ' + context.intervalValue);
		} else if (Number(context.intervalValue) == 0) {
			clearInterval(context.interval);
		}
	};

	context.loadData();
	context.interval = setInterval(context.loadData, Number(context.intervalValue));
});

app.directive('circle', `<div class="circle" ajsf-title="model.legend | suffix ' '| suffix model.label"><div class="circle-inner" ajsf-text="model.label"></div><span class="circle-legend" ajsf-text="model.legend"></span></div>`, (context, el) => {
	context.onRefresh = () => {
		if (context.model == undefined) {
			return;
		}

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

		$(el).find('.circle').attr('style', "--progress: " + context.model.percent + "; --fgColor: " + context.model.fgColor + "; --bgColor: " + context.model.bgColor + "; --size: " + size + "rem;");
	};
});

app.directive('bar', `<div class="bar" ajsf-title="model.legend | suffix ' '| suffix model.label"><span class="bar-legend" ajsf-text="model.legend"></span><div class="bar-inner"></div><span class="bar-label" ajsf-text="model.label"></span></div>`, (context, el) => {
	context.onRefresh = () => {
		if (context.model == undefined) {
			return;
		}

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