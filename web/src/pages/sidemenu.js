import React from "react"

const menuStyles = {
    fontSize: 30,
    listStyle: 'none',
    padding: 0, 
}

const menuElementStyle = {
    paddingBottom: 10,
    color: "#005618",
    cursor: 'pointer',
}

const menuActiveElementStyle = {
    paddingBottom: 10,
    color: "#14b341",
    cursor: 'pointer',
}

const SideMenu = ({ current, setCurrent }) => {
    return (
        <ul style={menuStyles}>
            <li style={current === "Home" ? menuActiveElementStyle : menuElementStyle} onClick={() => setCurrent("Home")}>
                <a>Home</a>
            </li>
            <li style={current === "Plan" ? menuActiveElementStyle : menuElementStyle} onClick={() => setCurrent("Plan")}>
                <a>Plan</a>
            </li>
            <li style={current === "Downloads" ? menuActiveElementStyle : menuElementStyle} onClick={() => setCurrent("Downloads")}>
                <a>Downloads</a>
            </li>
            <li style={current === "Contact" ? menuActiveElementStyle : menuElementStyle} onClick={() => setCurrent("Contact")}>
                <a>Contact</a>
            </li>
        </ul>
    )
}

export default SideMenu