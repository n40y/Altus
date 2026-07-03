// src/App.jsx
import React, { useEffect, useState, useRef } from 'react';
import CytoscapeComponent from 'react-cytoscapejs';
import axios from 'axios';
import './App.css';

export default function App() {
  const [elements, setElements] = useState([]);
  const [findings, setFindings] = useState([]);
  const [severityFilter, setSeverityFilter] = useState('ALL');
  const [loading, setLoading] = useState(true);
  
  const cyRef = useRef(null);

  useEffect(() => {
    async function initAltus() {
      try {
        const scanRes = await axios.get('http://localhost:8000/scan');
        setFindings(scanRes.data.findings);

        const graphRes = await axios.get('http://localhost:8000/graph');
        
        const cyElements = [];
        const processedNodes = new Set();

        graphRes.data.nodes.forEach(edge => {
          if (!processedNodes.has(edge.source)) {
            processedNodes.add(edge.source);
            cyElements.push({ data: { id: edge.source, type: edge.source_type } });
          }
          if (!processedNodes.has(edge.target)) {
            processedNodes.add(edge.target);
            cyElements.push({ data: { id: edge.target, type: edge.target_type } });
          }
          cyElements.push({ data: { source: edge.source, target: edge.target, label: edge.relation } });
        });

        setElements(cyElements);
      } catch (err) {
        console.error("Erreur de communication avec l'API Altus :", err);
      } finally {
        setLoading(false);
      }
    }
    initAltus();
  }, []);

  const exportAuditReport = () => {
    const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(findings, null, 2));
    const downloadAnchor = document.createElement('a');
    downloadAnchor.setAttribute("href", dataStr);
    downloadAnchor.setAttribute("download", "altus_audit_report.json");
    document.body.appendChild(downloadAnchor);
    downloadAnchor.click();
    downloadAnchor.remove();
  };

  const filteredFindings = findings.filter(f => severityFilter === 'ALL' || f.severity === severityFilter);

  const style = [
    // 1. RÈGLES GLOBALES
    {
      selector: 'edge[label = "CAN_ESCALATE_TO"]',
      style: {
        'line-color': '#ff9f43',        // Orange alerte
        'target-arrow-color': '#ff9f43',
        'width': 3,
        'line-style': 'dashed'          // Ligne pointillée pour symboliser le chemin d'attaque potentiel
      }
    },
    {
      selector: 'node',
      style: { 
        'color': '#fff', 
        'font-size': '11px', 
        'text-valign': 'center', 
        'text-halign': 'center',
        'width': '190px', 
        'height': '45px', 
        'shape': 'round-rectangle', 
        'border-width': 2 
      }
    },
    {
      selector: 'edge',
      style: { 
        'width': 2, 
        'line-color': '#4e5561', 
        'target-arrow-color': '#4e5561', 
        'target-arrow-shape': 'triangle', 
        'curve-style': 'bezier', 
        'label': 'data(label)', 
        'color': '#a1a8b5' 
      }
    },

    // 2. OVERRIDES SPÉCIFIQUES (Écrasent les règles globales)
    { 
      selector: 'node[type = "INTERNET"]', 
      style: { 
        'background-color': '#1e272e',   // Fond sombre pour éviter le conflit blanc sur blanc
        'border-color': '#ff4757',       // Bordure rouge vif (Alerte exposition)
        'border-width': '3px',
        'color': '#ff4757',               // Texte rouge assorti à la bordure
        'font-size': '12px',
        'font-weight': 'bold',
        'label': (ele) => `🌐 ${ele.data('id')}` // Utilise l'ID dynamique récupéré du scanner Go
      } 
    },
    { selector: 'node[type = "IAM_USER"]', style: { 'background-color': '#00adb5', 'border-color': '#007a80', 'label': (ele) => `👤 ${ele.data('id')}` } },
    { selector: 'node[type = "S3_BUCKET"]', style: { 'background-color': '#ff2e63', 'border-color': '#b8002f', 'label': (ele) => `🪣 ${ele.data('id')}` } },
    { selector: 'node[type = "IAM_POLICY"]', style: { 'background-color': '#f9ca24', 'border-color': '#c19a07', 'color': '#222831', 'label': (ele) => `📜 ${ele.data('id')}` } },
    { selector: 'node[type = "EC2_INSTANCE"]', style: { 'background-color': '#a55eea', 'border-color': '#7030b2', 'label': (ele) => `🖥️ ${ele.data('id')}` } }
  ];

  const handleCyInit = (cyInstance) => {
    cyRef.current = cyInstance;
    cyInstance.resize();
    cyInstance.fit();
  };

  if (loading) return <div style={{ color: '#fff', padding: '20px', backgroundColor: '#222831', minHeight: '100vh' }}>Analyse du graphe...</div>;

  return (
    <div style={{ backgroundColor: '#222831', minHeight: '100vh', padding: '24px', color: '#eee' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '20px' }}>
        <h1 style={{ margin: 0 }}>Altus - Cloud Risk Visualizer</h1>
        <div style={{ display: 'flex', gap: '10px' }}>
          <button onClick={() => setSeverityFilter('ALL')}>Tous</button>
          <button onClick={() => setSeverityFilter('CRITICAL')}>Critique</button>
          <button onClick={exportAuditReport}>📥 Exporter</button>
        </div>
      </div>

      <div style={{ display: 'flex', gap: '24px', height: 'calc(100vh - 100px)' }}>
        <div style={{ flex: 1, border: '1px solid #393e46', borderRadius: '8px', position: 'relative' }}>
          <CytoscapeComponent elements={elements} style={{ width: '100%', height: '100%' }} layout={{ name: 'cose' }} stylesheet={style} cy={handleCyInit} />
        </div>
        <div style={{ width: '360px', backgroundColor: '#393e46', padding: '20px', borderRadius: '8px' }}>
          <h2>Alertes</h2>
          {filteredFindings.map((f, i) => <div key={i}>[{f.severity}] {f.title}</div>)}
        </div>
      </div>
    </div>
  );
}