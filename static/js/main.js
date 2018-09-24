var mapName;
var map = null;

function getQueryVariable(variable) {
    var query = window.location.search.substring(1);
    var vars = query.split('&');
    for (var i = 0; i < vars.length; i++) {
        var pair = vars[i].split('=');
        if (decodeURIComponent(pair[0]) == variable) {
            return decodeURIComponent(pair[1]);
        }
    }
    console.log('Query variable %s not found', variable);
}

function getMap() {
    $.get("/api/map/" + mapName, function (data, status) {
        createMap(data);
    });
}

function newMap() {
    $.ajax({
        'type': 'POST',
        'url': "/api/map/",
        'contentType': 'application/json',
        'data': JSON.stringify({}),
        'dataType': 'json',
        'success': function (data, status) {
            window.location = "/?name=" + data.Name;
        }
    });
}


function addCountry(country) {
    $.ajax({
        'type': 'POST',
        'url': "/api/map/" + mapName + "/" + country + "/",
        'contentType': 'application/json',
        'data': {},
        'dataType': 'json',
        'success': function (data, status) {
            updateMap(data);
        }
    });
}

function updateMap(data) {
    calcData(data);
    map.updateChoropleth(data);
}

function calcData(data) {
    var onlyValues = Object.keys(data).map(function (key) {
        return data[key];
    });
    var maxValue = Math.max.apply(null, onlyValues);
    var paletteScale = d3.scale.linear()
        .domain([0, maxValue])
        .range(["#EFEFFF", "#02386F"]); // blue color
    for (const [key, value] of Object.entries(data)) {
        data[key] = {
            numberOfThings: value,
            fillColor: paletteScale(value)
        };
    }
}

function createMap(data) {
    console.log("create map");
    calcData(data);
    map = new Datamap({
        element: document.getElementById('container'),
        fills: {
            defaultFill: '#F5F5F5'
        },
        responsive: true,
        data: data,
        geographyConfig: {
            borderColor: '#DEDEDE',
            highlightBorderWidth: 2,
            // don't change color on mouse hover
            highlightFillColor: function (geo) {
                return geo['fillColor'] || '#F5F5F5';
            },
            // only change border
            highlightBorderColor: '#B7B7B7',
        },
        done: function (datamap) {
            datamap.svg.selectAll('.datamaps-subunit').on('click', function (geography) {
                addCountry(geography.id);
            });
        }
    });
    // Alternatively with jQuery
    $(window).on('resize', function () {
        map.resize();
    });
}

window.onload = function () {
    // Get mapName
    mapName = getQueryVariable("name");
    console.log(mapName);
    if (mapName == undefined) {
        newMap();
        return;
    }
    getMap();
}