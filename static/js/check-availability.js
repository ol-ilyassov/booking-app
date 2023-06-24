function checkAvailibility(room, CSRFToken) {
  document.getElementById("check-availability-btn").addEventListener("click", function () {
    let html = `
      <form class="container-fluid mt-1 mb-1" id="check-availability-form" action="#" method="post" novalidate class="needs-validation">
        <div class="row">
          <div class="col">
            <div class="row" id="reservation-dates-modal">
              <div class="col">
                <input disabled class="form-control" type="text" name="start_date" id="start_date" placeholder="Arrival" required>
              </div>
              <div class="col">
                <input disabled class="form-control" type="text" name="end_date" id="end_date" placeholder="Departure" required>
              </div>
            </div>
          </div>
        </div>
      </form>
    `
    attention.custom({
      msg: html,
      title: "Choose your dates",
      willOpen: () => {
        const elem = document.getElementById('reservation-dates-modal');
        const rp = new DateRangePicker(elem, {
          format: 'yyyy-mm-dd',
          showOnFocus: true,
          orientation: 'top',
          minDate: new Date(),
        });
      },
      preConfirm: () => {
        return [
          document.getElementById('start_date').value,
          document.getElementById('end_date').value
        ]
      },
      didOpen: () => {
        document.getElementById('start_date').disabled = false;
        document.getElementById('end_date').disabled = false;
      },
      callback: function (result) {
        let form = document.getElementById("check-availability-form");
        let formData = new FormData(form);
        formData.append("csrf_token", CSRFToken);
        formData.append("room_id", room);

        fetch('/search-availability-json', {
          method: "post",
          body: formData,
        })
          .then(response => response.json())
          .then(data => {
            if (data.ok) {
              attention.custom({
                icon: "success",
                showConfirmButton: false,
                msg: "<p>Room is available!</p>"
                  + "<p><a href='/book-room?id="
                  + data.room_id
                  + "&s="
                  + data.start_date
                  + "&e="
                  + data.end_date
                  + "' class='btn btn-primary'>"
                  + "Book Now!"
                  + "</a></p>",
              })
            } else {
              attention.error({
                msg: "Not available",
              })
            }
          })
      }
    });
  })
}