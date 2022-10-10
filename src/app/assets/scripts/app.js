// Returns navigator to be used by the game script in the iframe
var connectSerial = async function () {
  return navigator;
};

DOM.Ready(function () {
  // Insert CSRF tokens into forms
  window.onload = HandleDarkMode();

  // Perform AJAX post on click on method=post|delete anchors
  ActivateMethodLinks();

  // Insert CSRF tokens into forms
  ActivateForms();

  // Login
  ActivateCopyApprovedEmail();
  ActivateLoginInput();

  // Chat
  ActivateChart();
});

function ToggleDarkMode() {
  let bodyTag = document.getElementsByTagName("body")[0];
  let toggleTag = document.getElementById("colorToggle");

  if (bodyTag.classList.contains("lightMode")) {
    bodyTag.classList.replace("lightMode", "darkMode");
    toggleTag.innerHTML = "Light Mode";
  } else {
    bodyTag.classList.replace("darkMode", "lightMode");
    toggleTag.innerHTML = "Dark Mode";
  }
}

function HandleDarkMode() {
  let bodyTag = document.getElementsByTagName("body")[0];
  let toggleTag = document.getElementById("colorToggle");
  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    bodyTag.classList.add("darkMode");
    toggleTag.innerHTML = "Light Mode";
  } else {
    bodyTag.classList.add("lightMode");
    toggleTag.innerHTML = "Dark Mode";
  }

  toggleTag.addEventListener("click", ToggleDarkMode);
}

// Insert an input into every form with js to include the csrf token.
// this saves us having to insert tokens into every form.
function ActivateForms() {
  // Get authenticity token from head of page
  var token = authenticityToken();

  DOM.Each("form", function (f) {
    // Create an input element
    var csrf = document.createElement("input");
    csrf.setAttribute("name", "authenticity_token");
    csrf.setAttribute("value", token);
    csrf.setAttribute("type", "hidden");

    //Append the input
    f.appendChild(csrf);
  });
}

function authenticityToken() {
  // Collect the authenticity token from meta tags in header
  var meta = DOM.First("meta[name='authenticity_token']");
  if (meta === undefined) {
    e.preventDefault();
    return "";
  }
  return meta.getAttribute("content");
}

function ActivateCopyApprovedEmail() {
  DOM.On(".copy_button", "click", function (e) {
    /* Get the text field */
    var copyText = document.getElementById("approvedEmail");

    /* Select the text field */
    copyText.select();
    copyText.setSelectionRange(0, 99999); /* For mobile devices */

    /* Copy the text inside the text field */
    document.execCommand("copy");
  });
}

function ActivateLoginInput() {
  DOM.On(".code", "input", function (e) {
    var target = e.target,
      position = target.selectionEnd,
      length = target.value.length;

    target.value = target.value
      .replace(/[^\dA-Z]/g, "")
      .replace(/(.{4})/g, "$1 ")
      .trim();
    target.selectionEnd = position +=
      target.value.charAt(position - 1) === " " &&
      target.value.charAt(length - 1) === " " &&
      length !== target.value.length
        ? 1
        : 0;
  });
}

// Perform AJAX post on click on method=post|delete anchors
function ActivateMethodLinks() {
  DOM.On('a[method="post"]', "click", function (e) {
    var link = this;

    // Ignore disabled links
    if (DOM.HasClass(link, "disabled")) {
      e.preventDefault();
      return false;
    }

    // Get authenticity token from head of page
    var token = authenticityToken();

    // Perform a post to the specified url (href of link)
    var url = link.getAttribute("href");
    var data = "authenticity_token=" + token;

    DOM.Post(
      url,
      data,
      function (request) {
        // Use the response url to redirect
        window.location = request.responseURL;
      },
      function (request) {
        // Respond to error
        console.log("error", request);
      }
    );

    e.preventDefault();
    return false;
  });

  DOM.On('a[method="back"]', "click", function (e) {
    history.back(); // go back one step in history
    e.preventDefault();
    return false;
  });
}

// Handle activation and display of chart
function ActivateChart() {
  const CHART_COLORS = {
    red: "rgb(255, 99, 132)",
    orange: "rgb(255, 159, 64)",
    yellow: "rgb(255, 205, 86)",
    green: "rgb(75, 192, 192)",
    blue: "rgb(54, 162, 235)",
    purple: "rgb(153, 102, 255)",
    grey: "rgb(201, 203, 207)",
    maroon: "rgb(255,0,127)",
    pink: "rgb(255,0,255)",
  };

  const encodedCookieValue = getCookie("reports");

  if (encodedCookieValue.length > 0) {
    const decodedCookieValue = atob(encodedCookieValue);

    const reports = JSON.parse(decodedCookieValue);

    console.log("Value in reports after retrieving cookie: ", reports);

    let labels = [];
    let datasets = [];
    let apSpineTScore = [];
    let leftFemurNeckTScore = [];
    let rightFemurNeckTScore = [];
    let leftFemurTScore = [];
    let rightFemurTScore = [];
    let forearmTScore = [];
    let rightHipTScore = [];
    let leftHipTScore = [];

    for (let i = 0; i < reports.Dexas.length; i++) {
      var dexa = reports.Dexas[i];
      labels.push(dexa.Year);

      var organs = dexa.Organs;

      for (let j = 0; j < organs.length; j++) {
        var organ = organs[j];
        var site = organ.Direction + " " + organ.Site.replace("Total", "");
        site = site.trim().toLowerCase();

        switch (site) {
          case "ap spine l1-l4":
            apSpineTScore.push(organ.TScore);
            break;
          case "l1 through l4":
            apSpineTScore.push(organ.TScore);
            break;
          case "l1-l4":
            apSpineTScore.push(organ.TScore);
            break;
          case "ap lumbar spine":
            apSpineTScore.push(organ.TScore);
            break;                   
          case "left femur neck":
            leftFemurNeckTScore.push(organ.TScore);
            break;
          case "left femoral neck":
            leftFemurNeckTScore.push(organ.TScore);
            break;  
          case "right femur neck":
            rightFemurNeckTScore.push(organ.TScore);
            break;
          case "right femoral neck":
            rightFemurNeckTScore.push(organ.TScore);
            break;  
          case "left femur":
            leftFemurTScore.push(organ.TScore);
            break;
          case "right femur":
            rightFemurTScore.push(organ.TScore);
            break;
          case "forearm":
            forearmTScore.push(organ.TScore);
            break;
          case "left hip":
            leftHipTScore.push(organ.TScore);
            break;
          case "right hip":
            rightHipTScore.push(organ.TScore);  
        }
      }
    }

    console.log("Right femur Tscore: ", rightFemurTScore);

    if (apSpineTScore.length > 0) {
      datasets.push({
        label: "AP Spine L1-L4",
        borderColor: CHART_COLORS.red,
        backgroundColor: transparentize(CHART_COLORS.red, 0.5),
        data: apSpineTScore,
      });
    }

    if (leftFemurNeckTScore.length > 0) {
      datasets.push({
        label: "Left Femur Neck",
        borderColor: CHART_COLORS.blue,
        backgroundColor: transparentize(CHART_COLORS.blue, 0.5),
        data: leftFemurNeckTScore,
      });
    }

    if (rightFemurNeckTScore.length > 0) {
      datasets.push({
        label: "Right Femur Neck",
        borderColor: CHART_COLORS.orange,
        backgroundColor: transparentize(CHART_COLORS.orange, 0.5),
        data: rightFemurNeckTScore,
      });
    }

    if (leftFemurTScore.length > 0) {
      datasets.push({
        label: "Left Femur",
        borderColor: CHART_COLORS.green,
        backgroundColor: transparentize(CHART_COLORS.green, 0.5),
        data: leftFemurTScore,
      });
    }

    if (rightFemurTScore.length > 0) {
      datasets.push({
        label: "Right Femur",
        borderColor: CHART_COLORS.yellow,
        backgroundColor: transparentize(CHART_COLORS.yellow, 0.5),
        data: rightFemurTScore,
      });
    }

    if (forearmTScore.length > 0) {
      datasets.push({
        label: "Forearm",
        borderColor: CHART_COLORS.purple,
        backgroundColor: transparentize(CHART_COLORS.purple, 0.5),
        data: forearmTScore,
      });
    }

    if (rightHipTScore.length > 0) {
        datasets.push({
          label: "Right Hip",
          borderColor: CHART_COLORS.pink,
          backgroundColor: transparentize(CHART_COLORS.pink, 0.5),
          data: rightHipTScore,
        });
    }

    if (leftHipTScore.length > 0) {
        datasets.push({
          label: "Left Hip",
          borderColor: CHART_COLORS.maroon,
          backgroundColor: transparentize(CHART_COLORS.maroon, 0.5),
          data: leftHipTScore,
        });
    }

    if (reports.Dexas.length > 0) {
      ctx = document.getElementById("tScoreChart").getContext("2d");
      chart = new Chart(ctx, {
        type: "line",
        data: {
          labels: labels,
          datasets: datasets,
        },
        options: {
          responsive: true,
          maintainAspectRatio: true,
          plugins: {
            legend: {
              position: "top",
            },
            title: {
              display: true,
              text: "Bone Health Tracker",
            },
          },
        },
      });
    }
  }
}

function getCookie(cname) {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(";");
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == " ") {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

function transparentize(value, opacity) {
  var colorString = value,
    colorsOnly = colorString
      .substring(
        colorString.indexOf("(") + 1,

        colorString.lastIndexOf(")")
      )
      .split(/,\s*/),
    r = colorsOnly[0],
    g = colorsOnly[1],
    b = colorsOnly[2];

  const a = (1 - opacity) * 255;
  const calc = (x) => Math.round((x - a) / opacity);

  return `rgba(${calc(r)}, ${calc(g)}, ${calc(b)}, ${opacity})`;
}
