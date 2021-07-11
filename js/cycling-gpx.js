(function () {
    const trackStyle = { color: '#e600aa', weight: 4 }
    const distSpeedRange = 333

    let map = L.map('mapid')
    let trackLayer
    let velocityLayer

    Number.prototype.NN = function () { return (this < 10 ? '0' : '') + this }

    Date.prototype.date = function () {//в формат YYYY-MM-DD
        return this.getFullYear() + '-' + (this.getMonth() + 1).NN() + '-' + this.getDate().NN()
    }

    Date.prototype.time = function () {//в формат hh:mm:ss
        return this.getHours().NN() + ':' + this.getMinutes().NN() + ':' + this.getSeconds().NN()
    }

    Date.prototype.datetime = function () {//в формат YYYY-MM-DD hh:mm:ss
        return this.date() + ' ' + this.time()
    }

    function trackDist(track) {
        let dist = 0
        for (const m of track.dd) {
            dist += m
        }
        return dist
    }

    function initMap() {
        L.tileLayer('https://api.mapbox.com/styles/v1/{id}/tiles/{z}/{x}/{y}?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
            maxZoom: 18,
            attribution: '<a href="https://www.openstreetmap.org/copyright">© OpenStreetMap</a>',
            id: 'mapbox/streets-v11',
            tileSize: 512,
            zoomOffset: -1
        }).addTo(map);
    }

    function addTrack(track) {
        if (!!trackLayer) {
            trackLayer.remove()
        }
        trackLayer = L.polyline(track.ll, trackStyle)
        trackLayer.addTo(map)
        map.fitBounds(trackLayer.getBounds())
    }

    class SpeedFilter {
        constructor() {
            this.minSpeed = 10 //км/ч
            this.maxSpeed = 45 //км/ч
            this.maxAcceleration = 2 //км/ч/с

            this.prevTime = 0
            this.prevSpeed = 0
        }
        validate(time, speed) {
            // return true
            let acceleration = 0
            if (this.prevTime > 0) {
                let dt = time - this.prevTime
                if (dt > 0) {
                    acceleration = (speed - this.prevSpeed) / dt
                }
            }
            if (acceleration < this.maxAcceleration && speed < this.maxSpeed) {
                this.prevTime = time
                this.prevSpeed = speed
            }
            return acceleration < this.maxAcceleration && speed > this.minSpeed && speed < this.maxSpeed
        }
    }

    function addVelocities(timestamp, track) {
        if (!!velocityLayer) {
            velocityLayer.remove()
        }
        velocityLayer = L.layerGroup()
        let time = Number(timestamp)
        let dist = 0
        let totalDist = 0
        let markerDist = 0
        let markerTime = 0
        let maxSpeedIndex = 0
        let maxSpeed = 0
        let speedFilter = new SpeedFilter()
        for (const i in track.dd) {
            dist += track.dd[i]
            totalDist += track.dd[i]
            if (track.dt[i] > 0) {
                time += track.dt[i]
                let speed = track.dd[i] / track.dt[i] * 3.6 // скорость в точке трека в км/ч
                if (speedFilter.validate(time, speed)) {
                    if (speed > maxSpeed) {
                        maxSpeedIndex = i
                        maxSpeed = speed
                        markerDist = totalDist
                        markerTime = time
                    }
                }
            }
            if (dist > distSpeedRange) {
                dist = 0
                if (maxSpeedIndex > 0) {
                    title = `${maxSpeedIndex}\n`
                        + `${(new Date(markerTime * 1000)).time()}\n`
                        + `${(markerDist * 0.001).toFixed(1)} км`
                    L.marker(
                        track.ll[maxSpeedIndex],
                        {
                            icon: L.divIcon({
                                html: `<div title="${title}">${maxSpeed.toFixed(0)}</div>`,
                                iconSize: [30, 30]
                            })
                        }
                    ).addTo(velocityLayer)
                }
                maxSpeedIndex = 0
                maxSpeed = 0
            }
        }
        velocityLayer.addTo(map)
    }

    function initSelector() {
        let $select = $('select#tracks')
        let $option
        for (const ts in tracks) {
            let track = tracks[ts]
            let dist = trackDist(track) * 0.001
            let caption = `${(new Date(Number(ts) * 1000)).datetime()}&nbsp;&nbsp;&nbsp;${dist.toFixed(1)} км`
            $option = $(`<option value='${ts}'>${caption}</option>`)
            $select.append($option)
        }
        var hp = new URLSearchParams(location.hash.substr(1));
        if (hp.has('ts')) {// если параметра есть, то выделяем его, иначе выделяем последний
            $option = $select.find('option[value="'+hp.get('ts')+'"]')
        } else {
        }
        $option.prop('selected', true)
        $select.change(function () {
            let timestamp = $(this).val()
            let track = tracks[timestamp]
            addTrack(track)
            addVelocities(timestamp, track)
            location.hash = '#ts=' + timestamp
        })
        $select.trigger('change')
    }

    initMap();
    initSelector();
})()
