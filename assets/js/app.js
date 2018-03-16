$('#btn-test-get-info').click(function() {
  var serverIP = $('#serverIP').val();
  var statusMessage = $('#status-message');
  if (!serverIP) {
    statusMessage.html('Required serverIP field!').css('color', 'red');
    return;
  }
  $.get('/api/get-info/' + serverIP, function(data) {
    statusMessage
      .html('<pre>' + JSON.stringify(data, undefined, 2) + '</pre>')
      .css('color', 'blue');
  }).fail(function(data) {
    statusMessage.html(data.responseText).css('color', 'red');
  });
});

$('#btn-install').click(function() {
  var serverIP = $('#serverIP').val();
  var serverUsername = $('#serverUsername').val();
  var serverPassword = $('#serverPassword').val();
  var customer = $('#customer').val();
  var quotaTotal = $('#quotaTotal').val();
  var statusMessage = $('#status-message');
  if (
    !(serverIP && serverUsername && serverPassword && customer && quotaTotal)
  ) {
    statusMessage.html('Required all field!').css('color', 'red');
    return;
  }
  var btnInstall = $('#btn-install');
  statusMessage
    .html('Installing firmware, please wait...')
    .css('color', '#FFEB3B');
  btnInstall.attr('disabled', true);
  $.post(
    '/api/install',
    {
      serverIP: serverIP,
      serverUsername: serverUsername,
      serverPassword: serverPassword,
      customer: customer,
      quotaTotal: quotaTotal
    },
    function(data) {
      statusMessage.html(data).css('color', 'green');
      btnInstall.attr('disabled', false);
    }
  ).fail(function(data) {
    statusMessage.html(data.responseText).css('color', 'red');
    btnInstall.attr('disabled', false);
  });
});

function selectServer(ip) {
  $('#serverIP').val(ip);
}

$('#btn-scan').click(function() {
  var scanResult = $('#scan-result');
  var btnScan = $('#btn-scan');
  scanResult
    .html('<h3>Scanning IP, please wait...</h3>')
    .css('color', '#FFEB3B');
  btnScan.attr('disabled', true);
  $.get('/api/scan-ip', function(data) {
    if (!data) {
      scanResult.html('<h3>Not found</h3>').css('color', 'red');
    } else {
      var color = '';
      var html = '<ul class="server-list">';
      data.forEach(function(v) {
        if (v.firmwareInfo) {
          color = 'greenyellow';
        } else {
          color = '#9E9E9E';
        }
        html +=
          '<li style="background-color: ' +
          color +
          ';" onClick="selectServer(\'' +
          v.ip +
          '\')"><font color="blue">' +
          v.ip +
          '</font>';
        if (v.firmwareInfo) {
          html += '<ul>';
          html +=
            '<li>version: <font color="blue">' +
            v.firmwareInfo.version +
            '</font></li>';
          html +=
            '<li>buildTime: <font color="blue">' +
            v.firmwareInfo.buildTime +
            '</font></li>';
          html +=
            '<li>hardwareID: <font color="blue">' +
            v.firmwareInfo.hardwareID +
            '</font></li>';
          html +=
            '<li>customer: <font color="blue">' +
            v.firmwareInfo.customer +
            '</font></li>';
          html +=
            '<li>quotaTotal: <font color="blue">' +
            v.firmwareInfo.quota.total +
            '</font></li>';
          html +=
            '<li>quotaRemain: <font color="blue">' +
            v.firmwareInfo.quota.remain +
            '</font></li>';
          html += '</ul>';
        }
        html += '</li>';
      });
      html += '</ul>';
      scanResult.html(html);
    }
    btnScan.attr('disabled', false);
  }).fail(function(data) {
    scanResult.html(data.responseText).css('color', 'red');
    btnScan.attr('disabled', false);
  });
});
