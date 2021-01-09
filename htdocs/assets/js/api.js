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