import React from "react";
import '../style/downloads.css';

const downloadStyles = {
    borderRadius: 10,
    border: '4px solid #005618',
    padding: 15,
    height: '100%',
}

const Downloads = () => {
    const thesisURL = process.env.DOWNLOAD_THESIS_URL;
    const demoURL = process.env.DOWNLOAD_DEMO_URL;
    const executableURL = process.env.DOWNLOAD_EXECUTABLE_URL;
    const sourceCodeURL = process.env.DOWNLOAD_SOURCE_CODE_URL;
    return (
        <div style={downloadStyles} className={'downloads'}>
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