var remote = "http://192.168.99.119:8000";

function ready(fn) {
  if (document.readyState != 'loading') {
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

function get(selector) {
  return document.querySelectorAll(selector);
}

function create(eltype) {
  return document.createElement(eltype);
}

function foreach(elements, fn) {
  Array.prototype.forEach.call(elements, function(el, i) {
    fn(el, i);
  });
}

function parent(el) {
  return el.parentNode;
}

function append2Parent(pel, el) {
  pel.appendChild(el);
}

function onevent(el, eventName, eventHandler) {
  el.addEventListener(eventName, eventHandler);
}

function onclick(el, eventHandler) {
  onevent(el, "click", eventHandler);
}

function toggleClass(el, className) {
  el.classList.toggle(className);
}

function appendClass(el, className) {
  el.classList.add(className);
}

function settext(el, v) {
  el.textContent = v;
}

function request(method, endpoint, data, onloadHandler, onerrorHandler) {
  var request = new XMLHttpRequest();
  request.open(method, endpoint, true);
  request.onload = onloadHandler;
  request.onerror = onerrorHandler;
  if (data != null) {
    request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    request.send(data);
  } else {
    request.send();
  }
}

function appendDataAsTable(parent, endpoint, title) {

  request('GET', remote + endpoint, null,
    function() {
      var parel = get(parent)[0];

      var table = create("div");
      appendClass(table, "table");
      var header = create("div");
      appendClass(header, "table");
      appendClass(header, "header");
      settext(header, title);
      append2Parent(table, header);
      append2Parent(parel, table);
      if (this.status >= 200 && this.status < 400) {
        var resp = JSON.parse(this.response);
        if (resp.success) {
          if (resp.data.length > 0) {
            fi = resp.data[0];
            var thead = create("div");
            appendClass(thead, "table");
            appendClass(thead, "head");
            for (h in fi) {
              var c = create("div");
              appendClass(c, "table");
              appendClass(c, "head");
              appendClass(c, "cell");
              settext(c, h);
              append2Parent(thead, c);
            }
            var tbody = create("div");
            appendClass(tbody, "table");
            appendClass(tbody, "body");
            for (item in resp.data) {
              var r = create("div");
              appendClass(r, "table");
              appendClass(r, "body");
              appendClass(r, "row");
              for (h in resp.data[item]) {
                var c = create("div");
                appendClass(c, "table");
                appendClass(c, "head");
                appendClass(c, "cell");
                settext(c, resp.data[item][h]);
                append2Parent(r, c);
              }
              append2Parent(tbody, r);
            }
            append2Parent(table, thead);
            append2Parent(table, tbody);
          } else {
            var tbody = create(div);
            appendClass(tbody, "table");
            appendClass(tbody, "body");
            settext(tbody, "No " + title + " founded");
            append2Parent(table, tbody);
          }
        }
      } else {
        var errbody = create(div);
        appendClass(errbody, "table");
        appendClass(errbody, "error");
        settext(errbody, "Status: " + this.status + " Error: " + this.response);
        append2Parent(table, errbody);
      }
    },
    function() {
      var parel = get(parent)[0];
      var table = create("div");
      appendClass(table, "table");
      var header = create("div");
      appendClass(header, "table");
      appendClass(header, "header");
      settext(header, title);
      append2Parent(table, header);
      append2Parent(parel, table);
      errbody = create(div);
      appendClass(errbody, "table");
      appendClass(errbody, "error");
      settext(errbody, "Connection Error");
      append2Parent(table, errbody);
    }
  );
}

ready(function() {
  navlinks = get(".nav a");
  foreach(navlinks, function(el, i) {
    onclick(el, function(e) {
      e.preventDefault();
      toggleClass(get(".nav a .active")[0], "active");
      toggleClass(e.target, "active");
      panelid = parent(e.target).getAttribute('href');
      toggleClass(get(".panels .active")[0], "active");
      toggleClass(get(panelid)[0], "active");
    });
  });
});