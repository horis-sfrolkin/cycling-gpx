(function () {
    const trackStyle = { color: '#e600aa', weight: 4 }

    let map = L.map('mapid')
    let trackLayer

    Number.prototype.NN = function () { return (this < 10 ? '0' : '') + this }

    Date.prototype.date = function () {//в формат YYYY-MM-DD
        return this.getFullYear() + '-' + (this.getMonth() + 1).NN() + '-' + this.getDate().NN()
    }

    Date.prototype.datetime = function () {//в формат YYYY-MM-DD hh:mm:ss
        return this.date()
            + ' ' + this.getHours().NN() + ':' + this.getMinutes().NN() + ':' + this.getSeconds().NN()
    }

    function trackDist(track) {
        let dist = 0.0
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

    function initSelector() {
        let $select = $('select#tracks')
        let track
        let $option
        for (const ts in tracks) {
            track = tracks[ts]
            let dist = trackDist(track) * 0.001
            let caption = `${(new Date(Number(ts) * 1000)).datetime()}&nbsp;&nbsp;&nbsp;${dist.toFixed(3)} км`
            $option = $(`<option value='${ts}'>${caption}</option>`)
            $select.append($option)
        }
        $option.prop('selected', true)
        $select.change(function () {
            if (!!trackLayer) {
                trackLayer.remove()
            }
            trackLayer = L.polyline(tracks[$(this).val()].ll, trackStyle)
            trackLayer.addTo(map);
            map.fitBounds(trackLayer.getBounds())
        })
        $select.trigger('change')
    }

    initMap();
    initSelector();
})()
