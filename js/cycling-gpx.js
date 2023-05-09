(function () {
    const TRACK_STYLE = { color: '#e600aa', weight: 4 } // оформление трека
    const MIN_SPEED = 8 / 3.6 //минимальная скорость для учета активности м/с
    const MAX_SPPED = 60 / 3.6 //максимальная скорость для учета активности м/с
    const MAX_ACCELERATION = 4 / 3.6 //максимальное ускорение м/с2

    /** объекты карты */
    let map = L.map('mapid')
    /** Слой трека */
    let trackLayer
    /** Слой скоростей */
    let velocityLayer
    /** timestamp трека, который сейчас на карте*/
    let trackTime = 0
    /** данные построения последнего слоя скоростей, чтобы не строить его дважды */
    let velocityLayerTime = 0
    let velocityLayerZoom = 0

    Number.prototype.NN = function () { return (this < 10 ? '0' : '') + this }

    Number.prototype.time = function () {//в формат (?hh:)mm:ss
        let s = this % 60
        let m = Math.floor(this / 60) % 60
        let h = Math.floor(this / 3600)
        return (h < 1 ? '' : (h.NN() + ':')) + m.NN() + ':' + s.NN()
    }

    Date.prototype.date = function () {//в формат YYYY-MM-DD
        return this.getFullYear() + '-' + (this.getMonth() + 1).NN() + '-' + this.getDate().NN()
    }

    Date.prototype.time = function (onlyMinutes) {//в формат hh:mm:ss
        return this.getHours().NN()
            + ':' + this.getMinutes().NN()
            + (onlyMinutes ? '' : (':' + this.getSeconds().NN()))
    }

    Date.prototype.datetime = function () {//в формат YYYY-MM-DD hh:mm:ss
        return this.date() + ' ' + this.time(true)
    }

    function trackDist(track) {
        let dist = 0
        for (const m of track.dd) {
            dist += m
        }
        return dist
    }

    function initMap() {
        const url = 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
        const options = {
            attribution: '<a href="https://www.openstreetmap.org/copyright">© OpenStreetMap</a>'
        }
        L.tileLayer(url, options).addTo(map)
        map.on("zoomend", onZoomEnd)
    }

    function addTrack(track) {
        if (!!trackLayer) {
            trackLayer.remove()
        }
        trackLayer = L.polyline(track.ll, TRACK_STYLE)
        trackLayer.addTo(map)
        map.fitBounds(trackLayer.getBounds())
    }

    class SpeedFilter {
        constructor() {
            this.prevTime = 0
            this.prevSpeed = 0
        }
        validate(time, speed) {
            let acceleration = 0
            if (this.prevTime > 0) {
                let dt = time - this.prevTime
                if (dt > 0) {
                    acceleration = (speed - this.prevSpeed) / dt
                }
            }
            if (acceleration < MAX_ACCELERATION && speed < MAX_SPPED) {
                this.prevTime = time
                this.prevSpeed = speed
            }
            return acceleration < MAX_ACCELERATION && speed > MIN_SPEED && speed < MAX_SPPED
        }
    }

    function addVelocities(timestamp, zoom) {
        if (!!velocityLayer) {
            if (velocityLayerTime == timestamp && velocityLayerZoom == zoom) {
                return // этото слой скоростей уже построен - пропускаем
            }
            velocityLayer.remove()
        }
        velocityLayerTime = timestamp
        velocityLayerZoom = zoom
        velocityLayer = L.layerGroup()
        let distSpeedRange = 250 * Math.pow(2, (14 - zoom))
        let time = Number(timestamp)
        let dist = 0
        let totalDist = 0
        let speedFilter = new SpeedFilter()
        let track = tracks[timestamp]
        for (const i in track.dd) {
            dist += track.dd[i]
            totalDist += track.dd[i]
            if (track.dt[i] <= 0) {
                continue;
            }
            time += track.dt[i]
            let speed = track.dd[i] / track.dt[i]
            if (!speedFilter.validate(time, speed)) {
                continue;
            }
            if (dist >= distSpeedRange) {
                dist = 0
                title = `${new Date(time * 1000).time()} (${(time - timestamp).time()})`
                    + `\n${(totalDist * 0.001).toFixed(1)} км`
                L.marker(
                    track.ll[i],
                    {
                        icon: L.divIcon({
                            html: `<div title="${title}">${(3.6 * speed).toFixed(0)}</div>`,
                            iconSize: [30, 30]
                        })
                    }
                ).addTo(velocityLayer)
            }
        }
        velocityLayer.addTo(map)
    }

    function trackAverageSpeed(track) {
        let activeDist = 0
        let activeTime = 0
        let time = 0
        let speedFilter = new SpeedFilter()
        for (const i in track.dd) {
            if (track.dt[i] <= 0) {
                continue;
            }
            time += track.dt[i]
            let speed = track.dd[i] / track.dt[i]
            if (speedFilter.validate(time, speed)) {
                activeDist += track.dd[i]
                activeTime += track.dt[i]
            }
        }
        return activeTime > 0 ? (3.6 * activeDist / activeTime) : 0
    }

    function initSelector() {
        let $select = $('select#tracks')
        let $option
        for (const ts in tracks) {
            let track = tracks[ts]
            let dist = trackDist(track) * 0.001
            let speed = trackAverageSpeed(track)
            let caption = `${(new Date(Number(ts) * 1000)).datetime()}`
                + `&nbsp;&nbsp;&nbsp;${speed.toFixed(2)} км/ч`
                + `&nbsp;&nbsp;&nbsp;${dist.toFixed(1)} км`
            $option = $(`<option value='${ts}'>${caption}</option>`)
            $select.append($option)
        }
        var hp = new URLSearchParams(location.hash.substr(1));
        if (hp.has('ts')) {// если параметра есть, то выделяем его, иначе выделяем последний
            trackTime = Number(hp.get('ts'))
            $option = $select.find(`option[value="${trackTime}"]`)
        }
        $option.prop('selected', true)
        $select.change(function () {
            trackTime = Number($(this).val())
            let track = tracks[trackTime]
            addTrack(track)
            addVelocities(trackTime, map.getZoom())
            location.hash = `#ts=${trackTime}`
        })
        $select.trigger('change')
    }

    function onZoomEnd(ev) {
        addVelocities(trackTime, map.getZoom())
    }

    initMap();
    initSelector();
})()
