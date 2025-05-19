import { JSX } from 'solid-js';

interface SpinnerProps {
    show: boolean;
}

export const Spinner = (props: SpinnerProps): JSX.Element => {
    return (
        <span
        class='spinner'
        style={{
            'display': props.show ? 'inline-block' : 'none',
            'width': '20px',
            'height': '20px', 
            'border': '3px solid rgba(255, 255, 255, .3)', 
            'border-radius': '50%', 
            'border-top-color': '#fff', 
            'animation': 'spin 1s ease-in-out infinite', 
            'margin-left': '8px',
            'vertical-align': 'middle'
        }}
        />
    );
};