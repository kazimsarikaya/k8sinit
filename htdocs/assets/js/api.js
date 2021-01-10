function fillSummaryPanel() {
  appendDataAsTable("#summary", "/api/disks", "Block Devices");
  appendDataAsTable("#summary", "/api/zpools", "Zpools");
}

ready(fillSummaryPanel);

function sysaction(command, data) {
  request('POST', remote + "/api/system/" + command, data,
    function() {
      console.log(this.status, this.response);
    },
    function() {
      console.log("connection error");
    }
  );
}

ready(function() {
  var sysactions = get(".sysaction a");
  foreach(sysactions, function(el, i) {
    onclick(el, function(e) {
      e.preventDefault();
      command = parent(e.target).getAttribute('href').substring(1);
      var data = parent(e.target).getAttribute('action-data');
      sysaction(command, JSON.parse(data));
    });
  });
});

ready(function() {
  var nettypes = get(".nettype");
  foreach(nettypes, function(el, i) {
    onclick(el, function(e) {
      checked = e.target.checked;
      ipid = e.target.getAttribute('ip-id');
      var el = get("#" + ipid)[0];
      el.disabled = !checked;
      if (!checked) {
        el.value = "";
      }
      gwid = e.target.getAttribute('gw-id');
      if (gwid != null) {
        var el = get("#" + gwid)[0];
        el.disabled = !checked;
        if (!checked) {
          el.value = "";
        }
      }
    });
  });
});

ready(function() {
  request("GET", remote + "/api/disks", null,
    function() {
      if (this.status >= 200 && this.status < 400) {
        var resp = JSON.parse(this.response);
        var di = get("#disklist")[0];
        for (item in resp.data) {
          var o = create("option");
          settext(o, resp.data[item]["Path"]);
          append2Parent(di, o);
        }
      } else {
        console.log(this.status, this.response);
      }
    },
    function() {
      console.log("cannot get disk list");
    }
  );
  request("GET", remote + "/api/network/interfaces", null,
    function() {
      if (this.status >= 200 && this.status < 400) {
        var resp = JSON.parse(this.response);
        var iflists = get(".iflist");
        for (item in resp.data) {
          foreach(iflists, function(iflist, i) {
            var o = create("option");
            o.value = item;
            settext(o, item + " - " + resp.data[item]);
            append2Parent(iflist, o);
          });
        }
      } else {
        console.log(this.status, this.response);
      }
    },
    function() {
      console.log("cannot get disk list");
    }
  );
});

ready(function() {
  onclick(get(".installaction a")[0], function(e) {
    e.preventDefault();
    // Create WebSocket connection.
    var socket = new WebSocket('ws://192.168.99.119:8000/api/system/install');
    var installoutput = get("#installoutput")[0];

    // Connection opened
    socket.addEventListener('open', function(event) {
      socket.send('{ "disk": "/dev/sda", "force": true,   "poolname": "zp_k8s" }');
      console.log('data sended');
    });

    socket.addEventListener('close', function(event) {
      console.log('transaction ended');
    });

    // Listen for messages
    socket.addEventListener('message', function(event) {
      var line = create("div")
      settext(line, event.data);
      append2Parent(installoutput, line);
    });
  })
});