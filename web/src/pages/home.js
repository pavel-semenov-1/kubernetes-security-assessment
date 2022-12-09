import * as React from "react"
import '../style/home.css'

const homeStyles = {
    borderRadius: 10,
    border: '4px solid #005618',
    padding: 15,
    height: '100%',
    overflow: 'auto',
}

const Home = () => {
 return (
    <div style={homeStyles}>
        <h2 className="text-center">Home</h2>
        <p>Student Name: Pavel Semenov</p>
        <p>Study Programme: Applied Computer Science</p>
        <p>Primary language: English</p>
        <p>Secondary language: Slovak</p>
        <p>Supervisor: <a href="https://fmph.uniba.sk/en/contact/emplyoees/ostertag1">RNDr. Richard Ostertág, PhD.</a></p>
        <p>Abstract: Kubernetes has been gaining popularity rapidly in recent years as more and more enterprise solutions are subjected to cloud transformation and more companies are looking for the ways to increase development efficiency and reduce development costs. This brings new concerns from clients and stakeholders about the security of Kubernetes and its exposure to cyber-attacks.
        This thesis studies, compares and evaluates the state-of-the-art tools designed to discover vulnerabilities concerning the cluster configuration files, running pods or cluster itself. Assessment is carried out in both local cluster setup predisposed with multiple vulnerabilities and real-world enterprise cloud infrastructure. Based on the assessment results we intend either to improve one of the existing tools or develop a Kubernetes security framework of our own, which will be able to provide better results in addressing the cluster security.
        </p>
    </div>
 )   
}

export default Home