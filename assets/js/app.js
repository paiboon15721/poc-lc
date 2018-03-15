$('#btn-submit').click(function() {
  var serverIP = $('#serverIP').val();
  var serverUsername = $('#serverUsername').val();
  var serverPassword = $('#serverPassword').val();
  var customer = $('#customer').val();
  var quotaTotal = $('#quotaTotal').val();
  var message = $('#message');
  if (
    !(serverIP && serverUsername && serverPassword && customer && quotaTotal)
  ) {
    message.html('Required all field!').css('color', 'red');
    return;
  }
  var btnSubmit = $('#btn-submit');
  message.html('Installing firmware, please wait...').css('color', '#FFEB3B');
  btnSubmit.attr('disabled', true);
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
      message.html(data).css('color', '#64DD17');
      btnSubmit.attr('disabled', false);
    }
  ).fail(function(data) {
    message.html(data.responseText).css('color', 'red');
    btnSubmit.attr('disabled', false);
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
          color = 'red';
        }
        html +=
          '<li style="background-color: ' +
          color +
          ';" onClick="selectServer(\'' +
          v.ip +
          '\')">' +
          v.ip;
        if (v.firmwareInfo) {
          html += '<ul>';
          html += '<li>version: ' + v.firmwareInfo.version + '</li>';
          html += '<li>buildTime: ' + v.firmwareInfo.buildTime + '</li>';
          html += '<li>hardwareID: ' + v.firmwareInfo.hardwareID + '</li>';
          html += '<li>customer: ' + v.firmwareInfo.customer + '</li>';
          html += '<li>quotaTotal: ' + v.firmwareInfo.quota.total + '</li>';
          html += '<li>quotaRemain: ' + v.firmwareInfo.quota.remain + '</li>';
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
