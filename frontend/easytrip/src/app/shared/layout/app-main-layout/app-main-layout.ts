import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Header } from "../../header/header";

@Component({
  selector: 'app-app-main-layout',
  imports: [RouterOutlet, Header],
  templateUrl: './app-main-layout.html',
  styleUrl: './app-main-layout.css'
})
export class AppMainLayout {

}
