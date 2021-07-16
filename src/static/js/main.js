// Register button and input handler
window.onload = function() {
  var button = document.getElementById('scan-button');
  button.onclick = buttonHandler;

  // Pressing enter in the input box should simulate clicking the button
  var input = document.getElementById('host-input');
  input.addEventListener('keyup', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();
      button.click();
    }
  });
}

function buttonHandler() {
  // Get input
  var host = getHostInput();

  var button = document.getElementById('scan-button');
  if (host.length > 0 && !button.disabled) {
    button.value = 'Scanning';
    button.disabled = true;

    runNmapScan(host);
  }
}

function getHostInput() {
  return document.getElementById('host-input').value.trim();
}

function runNmapScan(host) {
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {
    if (this.readyState == 4) {
      // Don't forget to enable the button again!
      var button = document.getElementById('scan-button');
      button.value = 'Scan';
      button.disabled = false;

      if (this.status == 201) {
        var scan = JSON.parse(this.responseText);
        populateScanDiff(scan);
        populateScanHistory(host);
      } else if (this.status == 400) {
        displayHostError();
      } else {
        displayScanServerError();
      }
    }
  };
  xhr.open('POST', '/api/hosts/' + host + '/scans');
  xhr.send();
}

function populateScanDiff(scan) {
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        var previousScan = JSON.parse(this.responseText);
        displayDiff(scan, previousScan);
      } else if (this.status == 404) {
        displayNoPreviousScan(scan);
      } else {
        displayDiffServerError();
      }
    }
  };
  xhr.open('GET', '/api/scans/' + scan.Id + '/previousByHost');
  xhr.send();
}

function displayDiff(scan, previousScan) {
  var html = '';
  if (scan.Ports.length > 0) {
    html = '<p>Open ports: ' + scan.Ports.join(', ') + '</p>';
  } else {
    html = '<p>Open ports: None</p>';
  }

  var newlyOpened = scan.Ports.filter(x => !previousScan.Ports.includes(x));
  if (newlyOpened.length > 0) {
    html += '<p>Newly opened: ' + newlyOpened.join(', ') + '</p>';
  }

  var newlyClosed = previousScan.Ports.filter(x => !scan.Ports.includes(x));
  if (newlyClosed.length > 0) {
    html += '<p>Newly closed: ' + newlyClosed.join(', ') + '</p>';
  }

  document.getElementById('diff').innerHTML = html;
}

function displayNoPreviousScan(scan) {
  var html = '';
  if (scan.Ports.length > 1) {
    html = '<p>Open ports: ' + scan.Ports.join(', ') + '</p>';
  } else {
    html = '<p>Open ports: None</p>';
  }
  document.getElementById('diff').innerHTML = html;
}

function populateScanHistory(host) {
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        var scans = JSON.parse(this.responseText);
        displayHistory(host, scans);
      } else {
        displayHistoryServerError();
      }
    }
  };
  xhr.open('GET', '/api/hosts/' + host + '/scans');
  xhr.send();
}

function displayHistory(host, scans) {
  var html = '<div class="title">Scan history for ' + host + '</div>';
  for (scan of scans) {
    html += '<div class="scan">';
    html += '<p>Timestamp: ' + new Date(scan.DateTime).toLocaleString('en-US') + '</p>';
    if (scan.Ports.length > 0) {
      html += '<p>Ports: ' + scan.Ports.join(', ') + '</p>';
    } else {
      html += '<p>Ports: None</p>';
    }
    html += '</div>';
  }
  document.getElementById('history').innerHTML = html;
}

function displayHostError() {
  document.getElementById('history').innerHTML = '';
  document.getElementById('diff').innerHTML = '<p>Could not connect to the provided host. Please try again.</p>';
}

function displayScanServerError() {
  document.getElementById('history').innerHTML = '';
  document.getElementById('diff').innerHTML = '<p>There was a problem scanning for ports.</p>';
}

function displayDiffServerError() {
  document.getElementById('diff').innerHTML = '<p>There was a problem getting the diff.</p>';
}

function displayHistoryServerError() {
  document.getElementById('history').innerHTML = '<p>There was a problem getting the scan history.</p>';
}
