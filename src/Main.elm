module Main exposing (main)

import Base64
import Browser
import Bytes.Encode as Encode
import Flate
import Html exposing (Html)
import Html.Attributes as Attr
import Html.Events as Events
import Url



-- MAIN


main : Program () Model Msg
main =
    Browser.sandbox { init = init, update = update, view = view }



-- MODEL


type alias Model =
    { language : String
    , code : String
    }


init : Model
init =
    { language = "JavaScript"
    , code = "const foo = \"bar\""
    }



-- UPDATE


type Msg
    = TypedLanguage String
    | TypedCode String


update : Msg -> Model -> Model
update msg model =
    case msg of
        TypedLanguage l ->
            { model | language = l }

        TypedCode c ->
            { model | code = c }



-- VIEW


view : Model -> Html Msg
view model =
    Html.div [ Attr.class "relative w-screen h-screen monospace" ]
        [ Html.input
            [ Attr.class "fixed top-1/2 left-1/2"
            , Attr.class "p-1"
            , Attr.class "text-center"
            , Attr.class "text-white"
            , Attr.class "bg-black"
            , Attr.class "border-2 border-solid rounded-xl"
            , Attr.class "-translate-x-1/2 -translate-y-1/2"
            , Attr.title "Select language for syntax highlighting"
            , Attr.value model.language
            , Events.onInput TypedLanguage
            ]
            []
        , Html.div
            [ Attr.class "fixed right-0 bottom-0"
            , Attr.class "p-2"
            , Attr.class "text-center text-sm"
            ]
            [ Html.a [ Attr.href "https://enso.no", Attr.target "_blank" ] [ Html.text "Made with ❤️ by Ensō" ]
            , Html.text " | "
            , Html.a [ Attr.href "https://github.com/ensolabs/codimg", Attr.target "_blank" ] [ Html.text "Source" ]
            ]
        , Html.textarea
            [ halfScreenAttr
            , Attr.class "p-4"
            , Attr.class "font-mono"
            , Attr.class "border-2 border-solid rounded-xl"
            , Attr.class "outline-none"
            , Attr.class "resize-none"
            , Attr.value model.code
            , Events.onInput TypedCode
            , Attr.autofocus True
            ]
            []
        , Html.div [ halfScreenAttr, Attr.class "flex items-center justify-center" ]
            [ Html.img
                [ Attr.src <| codeToSrc model
                ]
                []
            ]
        ]


halfScreenAttr : Html.Attribute msg
halfScreenAttr =
    Attr.class "float-left w-full md:w-1/2 h-1/2 md:h-full overflow-y-scroll p-4"


codeToSrc : { a | code : String, language : String } -> String
codeToSrc { code, language } =
    --"https://codimg.alwaysdata.net/code.svg?input="
    "/code.svg?input="
        ++ encodeCodeBlock code
        ++ "&lang="
        ++ language


encodeCodeBlock : String -> String
encodeCodeBlock =
    Encode.string
        >> Encode.encode
        >> Flate.deflate
        >> Base64.fromBytes
        >> Maybe.withDefault "invalid data"
        >> Url.percentEncode
