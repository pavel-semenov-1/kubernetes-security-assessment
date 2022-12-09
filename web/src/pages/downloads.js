import React from "react";

const downloadStyles = {
    borderRadius: 10,
    border: '4px solid #005618',
    padding: 15,
    height: '100%',
}

const Downloads = () => {
    return (
        <div style={downloadStyles}>
            <h2 className="text-center">Downloads</h2>
            <ul>
                <li>
                    <a href="">Thesis</a>
                </li>
                <li>
                    <a href="">Demo</a>
                </li>
                <li>
                    <a href="">Executable</a>
                </li>
                <li>
                    <a href="">Source code</a>
                </li>
            </ul>
        </div>
    )
}

export default Downloads