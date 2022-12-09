import React from "react";
import "../style/contact.css"

const contactStyles = {
    borderRadius: 10,
    border: '4px solid #005618',
    padding: 15,
    height: '100%',
    overflow: 'auto',
}

const Contact = () => {
    return (
        <div style={contactStyles}>
            <h2 className="text-center">Contact</h2>
            <ul>
                <li>
                    <p>School email</p>
                    <a href="mailto:semenov1@uniba.sk">semenov1@uniba.sk</a>
                </li>
                <li>
                    <p>Work email</p>
                    <a href="mailto:Pavel.Semenov@ibm.com">Pavel.Semenov@ibm.com</a>
                </li>
                <li>
                    <p>Personal email</p>
                    <a href="mailto:q7.pavel@gmail.com">q7.pavel@gmail.com</a>
                </li>
                <li>
                    <p>LinkedIn</p>
                    <a href="https://www.linkedin.com/in/pavel-semenov-public/">Profile</a>
                </li>
            </ul>
        </div>
    )
}

export default Contact