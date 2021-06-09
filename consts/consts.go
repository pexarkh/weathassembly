package consts

const ForecastFragment = `
<header style="background-image: linear-gradient(rgba(0, 0, 0, 0.45), rgba(0, 0, 0, 0.75)), url({{ pfc.ImageUrl }})">
	<span>
		<h1>{{ pfc.PlaceName }}, <strong>{{ pfc.StateAbbreviation }}</strong></h1>
		<h2>{{ pfc | TopForecast }}</h2>
		<h3>
			<span>{{ pfc | TopTemperature }}</span> <sup>&deg;F</sup>
		</h3>
	</span>
	<span>{{ pfc | TopStartTime }}</span>
	<svg viewBox="0 0 100 17" preserveAspectRatio="none"><use href="#wave"></svg>
</header>
<section>
{% for fc in pfc.Forecasts %}
  <article>
  <h2>{{ fc.Name }}</h2>
  <p><svg><use href="#{{ fc.SvgSymbolId }}"></svg></p>
  <p>
    <span>{{ fc.Temperature }} &deg;F</span>
  </p>
  </article>
{% endfor %}
</section>`
