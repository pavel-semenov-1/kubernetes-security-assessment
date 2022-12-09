import React, { useState } from "react"
import "../style/index.css"
import "../fonts/CozetteVector.ttf"
import SideMenu from "./sidemenu"
import Home from "./home"
import Plan from "./plan"
import Downloads from "./downloads"
import Contact from "./contact"
import 'mdb-react-ui-kit/dist/css/mdb.min.css'
import "@fortawesome/fontawesome-free/css/all.min.css"

const pageStyles = {
  color: "#005618",
  fontFamily: "Cozette",
  flexDirection: 'row',
  height: '100%',
  padding: 96,
}
const sidebarStyles = {
  height: '100%',
}
const contentStyles = {
  fontSize: 24,
  marginLeft: 40,
  height: '100%',
  width: '100%',
}
const headingStyles = {
  marginTop: 0,
  marginBottom: 64,
  maxWidth: 300,
}
const headingAccentStyles = {
  color: "#134127",
}

const IndexPage = () => {
  const [current, setCurrent] = useState("Home");

  return (
    <main style={pageStyles} className={'d-flex flex-row'}>
      <div style={sidebarStyles}>
        <div id={"heading"} onClick={() => setCurrent("Home")} style={headingStyles}>
          <h1>
            Kubernetes Security Assessment
          </h1>
          <h3 style={headingAccentStyles}>
            Master Thesis
          </h3>
        </div>
        <SideMenu
          current={current}
          setCurrent={setCurrent}
        ></SideMenu>
      </div>
      <div style={contentStyles}>
        { 
          current === "Home" ? <Home></Home> :
          current === "Plan" ? <Plan></Plan> :
          current === "Downloads" ? <Downloads></Downloads> :
          current === "Contact" ? <Contact></Contact> : <p>Error</p>
        }
      </div>
    </main>
  )
}

export default IndexPage

export const Head = () => {
  const [mode, setMode] = useState("Dark");
  return (
    <>
      <title>Master thesis - Kubernetes security assessment</title>
      <style>{`body { background-color: ${mode === "Dark" ? "black" : "white"}; height: 100%; } html, #___gatsby, #gatsby-focus-wrapper { height: 100%; }`}</style>
    </>
  )
}