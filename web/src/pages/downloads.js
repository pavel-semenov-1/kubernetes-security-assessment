import React from "react";
import '../style/downloads.css';

const downloadStyles = {
    borderRadius: 10,
    border: '4px solid #005618',
    padding: 15,
    height: '100%',
}

const Downloads = () => {
    const thesisURL = process.env.DOWNLOAD_THESIS_URL || "https://pavel-semenov-1.github.io/kubernetes-security-assessment/artifacts/main.pdf";
    const demoURL = process.env.DOWNLOAD_DEMO_URL || "https://pavel-semenov-1.github.io/kubernetes-security-assessment";
    const seminarPresentationURL = process.env.DOWNLOAD_SEMINAR_PRESENTATION_URL || "https://pavel-semenov-1.github.io/kubernetes-security-assessment/artifacts/project_seminar.pptx";
    return (
        <div style={downloadStyles} className={'downloads'}>
            <h2 className="text-center">Downloads</h2>
            <ul>
                <li>
                    <a href={thesisURL}>Thesis</a>
                </li>
                <li>
                    <a href={demoURL}>Demo</a>
                </li>
                <li>
                    <a href={seminarPresentationURL}>Project Seminar 2 presentation</a>
                </li>
            </ul>
        </div>
    )
}

export default Downloads