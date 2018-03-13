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
