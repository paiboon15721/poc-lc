$('#btn-submit').click(function() {
  var serverIP = $('#serverIP').val();
  var serverUsername = $('#serverUsername').val();
  var serverPassword = $('#serverPassword').val();
  var customer = $('#customer').val();
  var quotaTotal = $('#quotaTotal').val();
  if (
    !(serverIP && serverUsername && serverPassword && customer && quotaTotal)
  ) {
    alert('Required all field!');
    return;
  }
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
      alert(data);
    }
  );
});
