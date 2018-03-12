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
  message.html('Installing firmware, please wait...').css('color', 'blue');
  btnSubmit.attr('disabled', true);
  $.post(
    '/install',
    {
      serverIP: serverIP,
      serverUsername: serverUsername,
      serverPassword: serverPassword,
      customer: customer,
      quotaTotal: quotaTotal
    },
    function(data) {
      message.html(data).css('color', 'green');
      btnSubmit.attr('disabled', false);
    }
  );
});
