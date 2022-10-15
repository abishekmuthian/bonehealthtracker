DOM.Ready(function () {
  // Handle theme
  window.onload = HandleDarkMode();

  // Perform AJAX post on click on method=post|delete anchors
  ActivateMethodLinks();

  // Insert CSRF tokens into forms
  ActivateForms();

  // Chart
  ActivateChart();
});

function ToggleDarkMode() {
  let bodyTag = document.getElementsByTagName("body")[0];
  let toggleTag = document.getElementById("colorToggle");

  if (bodyTag.classList.contains("lightMode")) {
    bodyTag.classList.replace("lightMode", "darkMode");
    toggleTag.innerHTML = "Light Mode";
    setCookie("theme", "dark");
  } else {
    bodyTag.classList.replace("darkMode", "lightMode");
    toggleTag.innerHTML = "Dark Mode";
    setCookie("theme", "light");
  }
}

function HandleDarkMode() {
  let toggleTag = document.getElementById("colorToggle");
  let bodyTag = document.getElementsByTagName("body")[0];
  // Not setting dark mode automatically due to the bug in Chromium on Linux.
  /*   
   if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    bodyTag.classList.add("darkMode");
    toggleTag.innerHTML = "Light Mode";
  } else {
    bodyTag.classList.add("lightMode");
    toggleTag.innerHTML = "Dark Mode";
  } */

  theme = getCookie("theme");

  if (theme) {
    if (theme === "light") {
      bodyTag.classList.add("lightMode");
      toggleTag.innerHTML = "Dark Mode";
    } else {
      bodyTag.classList.add("darkMode");
      toggleTag.innerHTML = "Light Mode";
    }
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
    var redirectURL = link.getAttribute("data-redirect");

    DOM.Post(
      url,
      data,
      function (request) {
        // Use the response url to redirect
        // window.location = request.responseURL;

        // Use the data attribute to redirect
        window.location = request.responseURL + redirectURL;
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

    // console.log("Value in reports after retrieving cookie: ", reports);

    let labels = [];
    let tScoreDatasets = [];
    let zScoreDatasets = [];
    let apSpineTScore = [];
    let apSpineZScore = [];
    let leftFemurNeckTScore = [];
    let leftFemurNeckZScore = [];
    let rightFemurNeckTScore = [];
    let rightFemurNeckZScore = [];
    let leftFemurTScore = [];
    let leftFemurZScore = [];
    let rightFemurTScore = [];
    let rightFemurZScore = [];
    let forearmTScore = [];
    let forearmZScore = [];
    let rightHipTScore = [];
    let rightHipZScore = [];
    let leftHipTScore = [];
    let leftHipZScore = [];

    for (let i = 0; i < reports.dexas.length; i++) {
      var dexa = reports.dexas[i];
      labels.push(dexa.year);

      var organs = dexa.organs;

      for (let j = 0; j < organs.length; j++) {
        var organ = organs[j];
        var site = organ.direction + " " + organ.site.replace("Total", "");
        site = site.trim().toLowerCase();

        switch (site) {
          case "ap spine l1-l4":
            apSpineTScore.push(organ.tScore);
            apSpineZScore.push(organ.zScore);
            break;
          case "l1 through l4":
            apSpineTScore.push(organ.tScore);
            apSpineZScore.push(organ.zScore);
            break;
          case "l1-l4":
            apSpineTScore.push(organ.tScore);
            apSpineZScore.push(organ.zScore);
            break;
          case "ap lumbar spine":
            apSpineTScore.push(organ.tScore);
            apSpineZScore.push(organ.zScore);
            break;
          case "left femur neck":
            leftFemurNeckTScore.push(organ.tScore);
            leftFemurNeckZScore.push(organ.zScore);
            break;
          case "left femoral neck":
            leftFemurNeckTScore.push(organ.tScore);
            leftFemurNeckZScore.push(organ.zScore);
            break;
          case "right femur neck":
            rightFemurNeckTScore.push(organ.tScore);
            rightFemurNeckZScore.push(organ.zScore);
            break;
          case "right femoral neck":
            rightFemurNeckTScore.push(organ.tScore);
            rightFemurNeckZScore.push(organ.zScore);
            break;
          case "left femur":
            leftFemurTScore.push(organ.tScore);
            leftFemurZScore.push(organ.zScore);
            break;
          case "right femur":
            rightFemurTScore.push(organ.tScore);
            rightFemurZScore.push(organ.zScore);
            break;
          case "forearm":
            forearmTScore.push(organ.tScore);
            forearmZScore.push(organ.zScore);
            break;
          case "left hip":
            leftHipTScore.push(organ.tScore);
            leftHipZScore.push(organ.zScore);
            break;
          case "right hip":
            rightHipTScore.push(organ.tScore);
            rightHipZScore.push(organ.zScore);
        }
      }
    }

    if (apSpineTScore.length > 0) {
      tScoreDatasets.push({
        label: "AP Spine L1-L4",
        borderColor: CHART_COLORS.red,
        backgroundColor: transparentize(CHART_COLORS.red, 0.5),
        data: apSpineTScore,
      });
    }

    if (apSpineZScore.length > 0) {
      zScoreDatasets.push({
        label: "AP Spine L1-L4",
        borderColor: CHART_COLORS.red,
        backgroundColor: transparentize(CHART_COLORS.red, 0.5),
        data: apSpineZScore,
      });
    }

    if (leftFemurNeckTScore.length > 0) {
      tScoreDatasets.push({
        label: "Left Femur Neck",
        borderColor: CHART_COLORS.blue,
        backgroundColor: transparentize(CHART_COLORS.blue, 0.5),
        data: leftFemurNeckTScore,
      });
    }

    if (leftFemurNeckZScore.length > 0) {
      zScoreDatasets.push({
        label: "Left Femur Neck",
        borderColor: CHART_COLORS.blue,
        backgroundColor: transparentize(CHART_COLORS.blue, 0.5),
        data: leftFemurNeckZScore,
      });
    }

    if (rightFemurNeckTScore.length > 0) {
      tScoreDatasets.push({
        label: "Right Femur Neck",
        borderColor: CHART_COLORS.orange,
        backgroundColor: transparentize(CHART_COLORS.orange, 0.5),
        data: rightFemurNeckTScore,
      });
    }

    if (rightFemurNeckZScore.length > 0) {
      zScoreDatasets.push({
        label: "Right Femur Neck",
        borderColor: CHART_COLORS.orange,
        backgroundColor: transparentize(CHART_COLORS.orange, 0.5),
        data: rightFemurNeckZScore,
      });
    }

    if (leftFemurTScore.length > 0) {
      tScoreDatasets.push({
        label: "Left Femur",
        borderColor: CHART_COLORS.green,
        backgroundColor: transparentize(CHART_COLORS.green, 0.5),
        data: leftFemurTScore,
      });
    }

    if (leftFemurZScore.length > 0) {
      zScoreDatasets.push({
        label: "Left Femur",
        borderColor: CHART_COLORS.green,
        backgroundColor: transparentize(CHART_COLORS.green, 0.5),
        data: leftFemurZScore,
      });
    }

    if (rightFemurTScore.length > 0) {
      tScoreDatasets.push({
        label: "Right Femur",
        borderColor: CHART_COLORS.yellow,
        backgroundColor: transparentize(CHART_COLORS.yellow, 0.5),
        data: rightFemurTScore,
      });
    }

    if (rightFemurZScore.length > 0) {
      zScoreDatasets.push({
        label: "Right Femur",
        borderColor: CHART_COLORS.yellow,
        backgroundColor: transparentize(CHART_COLORS.yellow, 0.5),
        data: rightFemurZScore,
      });
    }

    if (forearmTScore.length > 0) {
      tScoreDatasets.push({
        label: "Forearm",
        borderColor: CHART_COLORS.purple,
        backgroundColor: transparentize(CHART_COLORS.purple, 0.5),
        data: forearmTScore,
      });
    }

    if (forearmZScore.length > 0) {
      zScoreDatasets.push({
        label: "Forearm",
        borderColor: CHART_COLORS.purple,
        backgroundColor: transparentize(CHART_COLORS.purple, 0.5),
        data: forearmZScore,
      });
    }

    if (rightHipTScore.length > 0) {
      tScoreDatasets.push({
        label: "Right Hip",
        borderColor: CHART_COLORS.pink,
        backgroundColor: transparentize(CHART_COLORS.pink, 0.5),
        data: rightHipTScore,
      });
    }

    if (rightHipZScore.length > 0) {
      zScoreDatasets.push({
        label: "Right Hip",
        borderColor: CHART_COLORS.pink,
        backgroundColor: transparentize(CHART_COLORS.pink, 0.5),
        data: rightHipZScore,
      });
    }

    if (leftHipTScore.length > 0) {
      tScoreDatasets.push({
        label: "Left Hip",
        borderColor: CHART_COLORS.maroon,
        backgroundColor: transparentize(CHART_COLORS.maroon, 0.5),
        data: leftHipTScore,
      });
    }

    if (leftHipZScore.length > 0) {
      zScoreDatasets.push({
        label: "Left Hip",
        borderColor: CHART_COLORS.maroon,
        backgroundColor: transparentize(CHART_COLORS.maroon, 0.5),
        data: leftHipZScore,
      });
    }

    // console.log("datasets: ", datasets);

    if (reports.dexas.length > 0) {
      // Set T-score dataset
      tScoreCtx = document.getElementById("tScoreChart").getContext("2d");
      zScoreCtx = document.getElementById("zScoreChart").getContext("2d");

      tScoreChart = new Chart(tScoreCtx, {
        type: "line",
        data: {
          labels: labels,
          datasets: tScoreDatasets,
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
              text: "Bone Health Tracker (T-score)",
            },
          },
        },
      });

      zScoreChart = new Chart(zScoreCtx, {
        type: "line",
        data: {
          labels: labels,
          datasets: zScoreDatasets,
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
              text: "Bone Health Tracker (Z-score)",
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

function setCookie(cname, cvalue, exdays) {
  let expires;
  if (exdays) {
    const d = new Date();
    d.setTime(d.getTime() + exdays * 24 * 60 * 60 * 1000);
    expires = "expires=" + d.toUTCString();
  } else {
    expires = "Tue, 19 Jan 2038 04:14:07 GMT";
  }

  document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
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
