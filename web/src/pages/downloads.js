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
    const seminarPresentationURL = process.env.SEMINAR_PRESENTATION_URL;
    return (
        <div style={downloadStyles} className={'downloads'}>
            <h2 className="text-center">Downloads</h2>
            <ul>
                <li>
                    <a href={{thesisURL}}>Thesis</a>
                </li>
                <li>
                    <a href={{demoURL}}>Demo</a>
                </li>
                <li>
                    <a href={{seminarPresentationURL}}>Project Seminar 2 presentation</a>
                </li>
            </ul>
        </div>
    )
}

export default Downloads