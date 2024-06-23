import React, { useState } from 'react';
import './styles.css'; // Importa tu archivo CSS aquí

export function App() {
    const [arrayValue, setArrayValue] = useState('');
    const [result, setResult] = useState([]);

    const handleSubmit = async (event) => {
        event.preventDefault();

        try {
            const response = await fetch(`http://localhost:9000/predecir?arrayValue=${encodeURIComponent(arrayValue)}`);
            const data = await response.json();
            setResult(data);
        } catch (error) {
            console.error('Error fetching data:', error);
            setResult([]);
            alert('Error: No se pudo conectar con el servidor');
        }
    };

    return (
        <div className="container">
            <h1>Conversor de Temperatura</h1>
            <form onSubmit={handleSubmit}>
                <label htmlFor="arrayValue">Ingresar temperaturas en °C (Separar por espacios):</label>
                <input
                    type="text"
                    id="arrayValue"
                    name="arrayValue"
                    value={arrayValue}
                    onChange={(e) => setArrayValue(e.target.value)}
                />
                <button type="submit">Enviar</button>
            </form>

            {result.length > 0 && (
                <div>
                    <h2>Temperaturas en °F:</h2>
                    <ul>
                        {result.map((value, index) => (
                            <li key={index}>{value}</li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
}
